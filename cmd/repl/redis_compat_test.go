package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var replBinary string

func TestMain(m *testing.M) {
	if _, err := redisCommand("FLUSHDB"); err != nil {
		fmt.Fprintln(os.Stderr, "could not reset Redis before test run:", err)
		os.Exit(2)
	}

	buildDir, err := os.MkdirTemp("", "double-distributed-kv-tests-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create test build directory:", err)
		os.Exit(2)
	}
	replBinary = filepath.Join(buildDir, "repl")
	build := exec.Command("go", "build", "-o", replBinary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "could not build REPL: %v\n%s", err, output)
		_ = os.RemoveAll(buildDir)
		os.Exit(2)
	}

	code := m.Run()
	_ = os.RemoveAll(buildDir)
	os.Exit(code)
}

type compatibilityCase struct {
	command string
	input   string
	setup   []string
	redis   [][]string
	name    string
}

func TestRedisCommandCompatibility(t *testing.T) {
	cases := []compatibilityCase{
		{command: "SET", name: "new key", input: "SET ddkv:set:new value", redis: [][]string{{"SET", "ddkv:set:new", "value"}}},
		{command: "SET", name: "overwrite existing key", input: "SET ddkv:set:overwrite new", setup: []string{"SET ddkv:set:overwrite old"}, redis: [][]string{{"SET", "ddkv:set:overwrite", "old"}, {"SET", "ddkv:set:overwrite", "new"}}},
		{command: "SET", name: "empty value", input: `SET ddkv:set:empty ""`, redis: [][]string{{"SET", "ddkv:set:empty", ""}}},
		{command: "SET", name: "quoted key with spaces", input: `SET "ddkv set key" value`, redis: [][]string{{"SET", "ddkv set key", "value"}}},
		{command: "SET", name: "missing value argument", input: "SET ddkv:set:missing-value", redis: [][]string{{"SET", "ddkv:set:missing-value"}}},
		{command: "SET", name: "extra argument", input: "SET ddkv:set:extra value extra", redis: [][]string{{"SET", "ddkv:set:extra", "value", "extra"}}},

		{command: "GET", name: "existing key", input: "GET ddkv:get:existing", setup: []string{"SET ddkv:get:existing value"}, redis: [][]string{{"SET", "ddkv:get:existing", "value"}, {"GET", "ddkv:get:existing"}}},
		{command: "GET", name: "missing key", input: "GET ddkv:get:missing", redis: [][]string{{"GET", "ddkv:get:missing"}}},
		{command: "GET", name: "quoted key with spaces", input: `GET "ddkv get key"`, setup: []string{`SET "ddkv get key" value`}, redis: [][]string{{"SET", "ddkv get key", "value"}, {"GET", "ddkv get key"}}},
		{command: "GET", name: "missing key argument", input: "GET", redis: [][]string{{"GET"}}},
		{command: "GET", name: "extra argument", input: "GET ddkv:get:extra ignored", redis: [][]string{{"GET", "ddkv:get:extra", "ignored"}}},

		{command: "DEL", name: "delete existing key", input: "DEL ddkv:del:existing", setup: []string{"SET ddkv:del:existing value"}, redis: [][]string{{"SET", "ddkv:del:existing", "value"}, {"DEL", "ddkv:del:existing"}}},
		{command: "DEL", name: "delete missing key", input: "DEL ddkv:del:missing", redis: [][]string{{"DEL", "ddkv:del:missing"}}},
		{command: "DEL", name: "delete mix of existing and missing keys", input: "DEL ddkv:del:one ddkv:del:absent ddkv:del:two", setup: []string{"SET ddkv:del:one 1", "SET ddkv:del:two 2"}, redis: [][]string{{"SET", "ddkv:del:one", "1"}, {"SET", "ddkv:del:two", "2"}, {"DEL", "ddkv:del:one", "ddkv:del:absent", "ddkv:del:two"}}},
		{command: "DEL", name: "duplicate key", input: "DEL ddkv:del:duplicate ddkv:del:duplicate", setup: []string{"SET ddkv:del:duplicate value"}, redis: [][]string{{"SET", "ddkv:del:duplicate", "value"}, {"DEL", "ddkv:del:duplicate", "ddkv:del:duplicate"}}},
		{command: "DEL", name: "no key arguments", input: "DEL", redis: [][]string{{"DEL"}}},
		{command: "DEL", name: "quoted key with spaces", input: `DEL "ddkv del key"`, setup: []string{`SET "ddkv del key" value`}, redis: [][]string{{"SET", "ddkv del key", "value"}, {"DEL", "ddkv del key"}}},

		{command: "EXISTS", name: "existing key", input: "EXISTS ddkv:exists:existing", setup: []string{"SET ddkv:exists:existing value"}, redis: [][]string{{"SET", "ddkv:exists:existing", "value"}, {"EXISTS", "ddkv:exists:existing"}}},
		{command: "EXISTS", name: "missing key", input: "EXISTS ddkv:exists:missing", redis: [][]string{{"EXISTS", "ddkv:exists:missing"}}},
		{command: "EXISTS", name: "mix of existing and missing keys", input: "EXISTS ddkv:exists:one ddkv:exists:absent ddkv:exists:two", setup: []string{"SET ddkv:exists:one 1", "SET ddkv:exists:two 2"}, redis: [][]string{{"SET", "ddkv:exists:one", "1"}, {"SET", "ddkv:exists:two", "2"}, {"EXISTS", "ddkv:exists:one", "ddkv:exists:absent", "ddkv:exists:two"}}},
		{command: "EXISTS", name: "duplicate key", input: "EXISTS ddkv:exists:duplicate ddkv:exists:duplicate", setup: []string{"SET ddkv:exists:duplicate value"}, redis: [][]string{{"SET", "ddkv:exists:duplicate", "value"}, {"EXISTS", "ddkv:exists:duplicate", "ddkv:exists:duplicate"}}},
		{command: "EXISTS", name: "no key arguments", input: "EXISTS", redis: [][]string{{"EXISTS"}}},
		{command: "EXISTS", name: "quoted key with spaces", input: `EXISTS "ddkv exists key"`, setup: []string{`SET "ddkv exists key" value`}, redis: [][]string{{"SET", "ddkv exists key", "value"}, {"EXISTS", "ddkv exists key"}}},
	}

	passed := 0
	failed := 0
	passedCases := make([]compatibilityCase, 0, len(cases))
	for _, tc := range cases {
		t.Run(tc.command+"/"+tc.name, func(t *testing.T) {
			commands := append(append([]string{}, tc.setup...), tc.input)
			appOutput, err := runREPL(strings.Join(commands, "\n") + "\n")
			if err != nil {
				t.Fatalf("could not run REPL: %v", err)
			}

			var expected string
			for _, command := range tc.redis {
				expected, err = redisCommand(command...)
				if err != nil {
					t.Fatalf("could not run redis-cli %q: %v", command, err)
				}
			}
			actual := lastREPLResponse(appOutput)
			if actual != expected {
				failed++
				t.Errorf("command: %s\nexpected (redis-cli): %q\nactual (REPL):       %q", tc.input, expected, actual)
			} else {
				passed++
				passedCases = append(passedCases, tc)
			}
		})
	}
	for _, tc := range passedCases {
		t.Logf("PASS: %s | %s", tc.name, tc.input)
	}
	t.Logf("STATS: total=%d passed=%d failed=%d", len(cases), passed, failed)
}

func runREPL(input string) (string, error) {
	command := exec.Command(replBinary)
	command.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stdout
	err := command.Run()
	return stdout.String(), err
}

func lastREPLResponse(output string) string {
	lines := strings.Split(strings.ReplaceAll(output, "double-distributed-kv> ", ""), "\n")
	for _, line := range slices.Backward(lines) {
		if response := strings.TrimSpace(line); response != "" {
			return response
		}
	}
	return ""
}

func redisCommand(args ...string) (string, error) {
	commandArgs := append([]string{"--raw"}, args...)
	command := exec.Command("redis-cli", commandArgs...)
	output, err := command.CombinedOutput()
	return strings.TrimRight(string(output), "\r\n"), err
}

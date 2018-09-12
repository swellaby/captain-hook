// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	testResultsDirectory = ".testresults/"
	coverageResultsDirectory = ".coverage/"
	junitXmlTestResultsFileName = "junit.xml"
	junitXmlTestResultsFile = testResultsDirectory + junitXmlTestResultsFileName
	jsonTestResultsFileName = "unit.xml"
	jsonTestResultsFile = testResultsDirectory + jsonTestResultsFileName
	goTestResultsJsonFileName = "unit.json"
	coverageOutFileName = "coverage.out"
	coberturaCoverageFileName = "cobertura.xml"
	coverageOutFile = coverageResultsDirectory + coverageOutFileName
	coberturaCoverageFile = coverageResultsDirectory + coberturaCoverageFileName
)

func goGetTool(tool string) {
	cmd := exec.Command("go", "get", tool)
	cmd.Dir = os.TempDir();
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", string(out))

	if err != nil {
		os.Exit(1)
	}
}

func installDevTools() {
	fmt.Println("Installing dev tools...")
	fmt.Println("Installing go-junit-report...")
	goGetTool("github.com/jstemmer/go-junit-report")
}

func createTestOutputDirectories() {
	os.MkdirAll(testResultsDirectory, os.ModePerm)
	os.MkdirAll(coverageResultsDirectory, os.ModePerm)
}

func cleanTestResultFiles() {
	os.RemoveAll(testResultsDirectory)
	os.RemoveAll(coverageResultsDirectory)
}

func Setup() {
	installDevTools()
	fmt.Println("Creating directories for test/coverage/etc. output files...")
	createTestOutputDirectories()
}

func getPackageNames(pkgPath string) string {
	cmd := exec.Command("go", "list", pkgPath)
	out, _ := cmd.CombinedOutput()
	result := strings.TrimSuffix(string(out), "\n")
	return result

	// packages := strings.Split(string(out), "\n")
	// for _, pkg := range packages {
	// 	if len(pkg) < 1 {
	// 		fmt.Println("nope, last one")
	// 	} else {
	// 		fmt.Printf("%s\n", pkg)
	// 	}
	// }
}

// Test Runs the unit tests
func Test() error {
	cleanTestResultFiles()
	createTestOutputDirectories()
	fmt.Println("Running tests...")
	pkg := getPackageNames("./pkg/...")
	coverProfile := "-coverprofile=" + coverageOutFile
	testCmd := exec.Command("go", "test", pkg, "-v", coverProfile)
	testOutput, err := testCmd.CombinedOutput()
	fmt.Printf(string(testOutput))

	// Create Standard JSON Result File
	jsonCmd := exec.Command("go", "tool", "test2json", "-t", "-p", pkg)
	jsonCmd.Stdin = bytes.NewBuffer(testOutput)
	jsonOutFile, err := os.Create(jsonTestResultsFile)
	defer jsonOutFile.Close()
	jsonCmd.Stdout = jsonOutFile
	jsonCmd.Run()


	// Create JUnit XML Result File
	junitCmd := exec.Command("go-junit-report")
	junitCmd.Stdin = bytes.NewBuffer(testOutput)
	outfile, err := os.Create(junitXmlTestResultsFile)
	defer outfile.Close()
	junitCmd.Stdout = outfile
	junitCmd.Run()

	return err
}

// Lint Runs the linter
func Lint() error {
	fmt.Println("Running golint...")
	cmd := exec.Command("golint", "./...")
	output, err := cmd.CombinedOutput()
	fmt.Printf(string(output))
	return err
}

// A build step that requires additional params, or platform specific steps for example
// func Build() error {
// 	fmt.Println("Building...")
// 	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
// 	return cmd.Run()
// }

func Ci() {
	Test()
	fmt.Println()
	Lint()
}

func Clean() {
	fmt.Println("Cleaning...")
	//os.RemoveAll("MyApp")
}

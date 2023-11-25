# MutableIn

MutableIn is a Go library that provides an os.Stdin-like stream for custom interactive input, enabling external modification and real-time reflection in the CLI.

# Table Of Contents
- [Installation](#installation)
- [Usage](#Usage)
- [Testing and Considerations](#testing-and-considerations)

# Installation
To install MutableIn, use the following command:
```bash
go get github.com/avearmin/mutable-in
```

# Usage
To incorporate MutableIn into your Go project, follow these steps:

1. Initialize MutableIn:
```Go
muIn := mutablein.NewMutableIn()
muIn.Init()
defer muIn.Close()
```
2. Pass MutableIn to a scanner:
```Go
scn := bufio.NewScanner(muIn)
```
That's it! Now, MutableIn will capture input from the keyboard like os.Stdin. The key distinction lies in the ability to write into it from your codebase:
```Go
muIn.Write([]byte("Hello world!"))
```
This input will seamlessly appear in your terminal, just as if it were typed from the keyboard. It can be added to with any key, and erased with backspace. This unique feature enhances the flexibility of interactive input handling in your CLI applications.

# Testing and Considerations

**Warning: MutableIn is in Active Development**

MutableIn is currently in active development, and at its current stage, it should be considered a proof of concept. The primary focus of the project is on closely mimicking the visual behavior of your terminal. Due to this emphasis, testing is crucial to ensure the desired functionality.

## Manual Testing Required
As of now, much of MutableIn's functionality relies on visual conformity with your terminal. Therefore, extensive manual testing is necessary to validate its behavior in different terminal environments. Keep in mind that MutableIn may not work seamlessly across all terminals, as different terminals can represent keys using distinct byte sequences. It is essential to perform manual tests to verify compatibility.

## Terminal-Specific Byte Representations
MutableIn's behavior relies on accurate recognition of key sequences, and variations in how terminals represent keys can impact its functionality. Testing for these representations must be done manually to ensure proper operation. While MutableIn is expected to work with most Unix terminals, comprehensive testing is still required to offer full confidence in its compatibility.

## Compatibility Considerations
As of the current development stage, MutableIn is not compatible with Windows Command Prompt or PowerShell. If you are using these environments, please be aware that MutableIn may not function as intended.
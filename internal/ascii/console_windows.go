//go:build windows
// +build windows

package ascii

import (
    "syscall"
    "unsafe"
    "os"
)

var (
    kernel32                       = syscall.NewLazyDLL("kernel32.dll")
    procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
    procGetConsoleCursorInfo       = kernel32.NewProc("GetConsoleCursorInfo")
    procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
    procSetConsoleCursorInfo       = kernel32.NewProc("SetConsoleCursorInfo")
    procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
)

type short int16
type dword uint32
type word uint16

type coord struct {
    x short
    y short
}

type smallRect struct {
    bottom short
    left   short
    right  short
    top    short
}

type consoleScreenBufferInfo struct {
    size              coord
    cursorPosition    coord
    attributes        word
    window            smallRect
    maximumWindowSize coord
}

type consoleCursorInfo struct {
    size    dword
    visible int32
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func ShowCursor() {
    handle := syscall.Handle(os.Stderr.Fd())

    var cci consoleCursorInfo
    _, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
    cci.visible = 1

    _, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func HideCursor() {
    handle := syscall.Handle(os.Stderr.Fd())

    var cci consoleCursorInfo
    _, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
    cci.visible = 0

    _, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// ClearLine clears the current line and moves the cursor to its start position.
func ClearLine() {
    handle := syscall.Handle(os.Stderr.Fd())

    var csbi consoleScreenBufferInfo
    _, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

    var w uint32
    var x short
    cursor := csbi.cursorPosition
    x = csbi.size.x
    _, _, _ = procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
}

// Clear clears the current position and moves the cursor to the left.
func Clear() {
    handle := syscall.Handle(os.Stderr.Fd())

    var csbi consoleScreenBufferInfo
    _, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

    var w uint32
    cursor := csbi.cursorPosition
    _, _, _ = procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(1), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))

    if cursor.x > 0 {
        cursor.x--
    }
    _, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}
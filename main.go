package main
import (
    "log"
    "pku/mpk"
    "unsafe"
    "fmt"
)

func main() {
// Allocate an array, large enough to likely be on its own page
    a := make([]int, 1, 10000)

    // Allocate a protection key
    pkey, err := mpk.PkeyAlloc()
    if err != nil {
        log.Fatalf("Failed to allocate protection key: %v", err)
    }
    defer mpk.PkeyFree(pkey)  // Ensure the key is freed at the end

    // Align pointer to page boundary and protect the page as read-only
    alignedAddr := (uintptr(unsafe.Pointer(&a[0])) >> 12) << 12
    err = mpk.PkeyMprotect(
        alignedAddr,
        1<<12,          // Page size (typically 4096 bytes or 4KiB)
        mpk.SysProtR, // Set to read-only
        pkey,           // Use the allocated key
    )
    if err != nil {
        log.Fatalf("Failed to set memory protection: %v", err)
    }

    // Attempt to write to the protected memory; this should fail
    fmt.Printf("Attempting to write to read-only memory...\n")
    a[0] = 123  // This line should cause a segmentation fault if protection is enforced

    // If the program continues, print the value
    fmt.Printf("Value at a[0]: %d\n", a[0])
}


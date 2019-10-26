## Build

```console
$ go build
```

## Usage

```console
Usage: simulator [OPTIONS] FILENAME ENTRYPOINT
  -debug
        write debug log to stderr
```

## Test

```console
$ go test ./pkg/simulator
```

## Example Usage

```
$ cat > fib.s < EOF
.data
.text
fib.9:
addi $c1, $zero, 2
slt $c1, $i0, $c1
beq $c1, $zero, 1
j ifge_else.23
addi $cl, $i0, -1
sw $i0, 0($sp)
add $i0, $zero, $cl
sw $ra, 4($sp)
addi $sp, $sp, 8
jal fib.9
addi $sp, $sp, -8
lw $ra, 4($sp)
add $cl, $zero, $i0
lw $i0, 0($sp)
addi $i0, $i0, -2
sw $cl, 4($sp)
sw $ra, 12($sp)
addi $sp, $sp, 16
jal fib.9
addi $sp, $sp, -16
lw $ra, 12($sp)
add $cl, $zero, $i0
lw $i0, 4($sp)
add $i0, $i0, $cl
jr $ra
ifge_else.23:
addi $i0, $zero, 1
jr $ra
start:
addi $hp $zero 1000
addi $sp $zero 2000
addi $i0, $zero, 10
sw $ra, 4($sp)
addi $sp, $sp, 8
jal fib.9
addi $sp, $sp, -8
lw $ra, 4($sp)
sw $ra, 4($sp)
addi $sp, $sp, 8
jal print_int
addi $sp, $sp, -8
lw $ra, 4($sp)
exit
print_int:
out $i0
jr $ra
EOF
$ ./simulator ./fib.s start
89
```

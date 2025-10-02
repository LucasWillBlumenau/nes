from dataclasses import dataclass


@dataclass(slots=True)
class Instruction:
    name: str
    address: int
    x: int
    y: int
    a: int
    sp: int
    p: int 

    def __eq__(self, other) -> None:
        return (
            isinstance(other, Instruction)
            and self.name == other.name
            and self.address == other.address
            and self.x == other.x
            and self.y == other.y
            and self.a == other.a
            and self.sp == other.sp
            and self.p == other.p
        )



def main() -> None:
    output_lines = read_file("output.log").split('\n')
    nestest_lines = read_file("nestest.log").split('\n')
    mix_size = min([len(output_lines), len(nestest_lines)])
    for i in range(mix_size):
        output = parse_output_line(output_lines[i].strip())
        nestest = parse_nestest_line(nestest_lines[i].strip())

        if output != nestest:
            print(f"Problem in line {i + 1}")



def read_file(path: str) -> str:
    with open(path, encoding='UTF8') as fp:
        return fp.read()


def parse_nestest_line(line: str) -> Instruction:  
    parts = line.split()
    i = 0
    
    address = int(parts[i], 16)
    i += 1

    while len(parts[i]) == 2:
        i += 1
    
    name = parts[i].upper().strip()
    i += 1

    while parts[i].startswith('A:') is False:
        i += 1

    a = int(parts[i].removeprefix('A:'), 16)
    i += 1

    x = int(parts[i].removeprefix('X:'), 16)
    i += 1

    y = int(parts[i].removeprefix('Y:'), 16)
    i += 1

    p = int(parts[i].removeprefix('P:'), 16)
    i += 1

    sp = int(parts[i].removeprefix('SP:'), 16)
    i += 1

    return Instruction(
        name=name,
        address=address,
        x=x,
        y=y,
        a=a,
        sp=sp,
        p=p,
    )


def is_hex(value: str) -> bool:
    for c in value:
        if is_hex_char(c) is False:
            return False
    return True


def is_hex_char(char: str) -> bool:
    return (
        (char >= '0' and char <= '9')
        or (char >= 'A' and char <= 'F')
        or (char >= 'a' and char <= 'f')
    )


def parse_output_line(line: str) -> Instruction:
    parts = [p.strip() for p in line.split()]
    i = 0
    address = int(parts[i], 16)
    i += 1

    name = parts[i].upper().strip()
    i += 1

    while parts[i].startswith("A:") is False:
        i += 1
    i += 1

    a = int(parts[i].removeprefix("$").removesuffix(','), 16)
    i += 2
    x = int(parts[i].removeprefix("$").removesuffix(','), 16)
    i += 2
    y = int(parts[i].removeprefix("$").removesuffix(','), 16)
    i += 2
    p = int(parts[i].removeprefix("$").removesuffix(','), 16)
    i += 2
    sp = int(parts[i].removeprefix("$01").removesuffix(','), 16)

    return Instruction(
        name=name,
        address=address,
        a=a,
        x=x,
        y=y,
        p=p,
        sp=sp,
    )


if __name__ == '__main__':
    main()
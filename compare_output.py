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
    clock: int

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
            and self.clock == other.clock
        )



def main() -> None:
    output_lines = read_file("output.log").split('\n')
    nestest_lines = read_file("nestest.log").split('\n')
    mix_size = min([len(output_lines), len(nestest_lines)])
    for i in range(mix_size):
        output = parse_output_line(output_lines[i].strip())
        nestest = parse_nestest_line(nestest_lines[i].strip())
        if output == nestest:
            continue

        print(f"Problem in line {i + 1}")
        if output.name != nestest.name:
            print(f"output name: {output.name}, nestest name: {nestest.name} ")
        if output.address != nestest.address:
            print(f"output address: {output.address}, nestest address: {nestest.address} ")
        if output.x != nestest.x:
            print(f"output x: {output.x}, nestest x: {nestest.x} ")
        if output.y != nestest.y:
            print(f"output y: {output.y}, nestest y: {nestest.y} ")
        if output.a != nestest.a:
            print(f"output a: {output.a}, nestest a: {nestest.a} ")
        if output.sp != nestest.sp:
            print(f"output sp: {output.sp}, nestest sp: {nestest.sp} ")
        if output.p != nestest.p:
            print(f"output p: {output.p}, nestest p: {nestest.p} ")
        if output.clock != nestest.clock:
            print(f"output clock: {output.clock}, nestest clock: {nestest.clock} ")
        break



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

    clock = int(parts[i:][-1].removeprefix('CYC:'))

    return Instruction(
        name=name,
        address=address,
        x=x,
        y=y,
        a=a,
        sp=sp,
        p=p,
        clock=clock,
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
    name = parts[i].upper().strip()
    i += 1

    while parts[i].startswith("A:") is False:
        i += 1
    i += 1

    a = int(parts[i].removesuffix(','), 16)
    i += 2
    x = int(parts[i].removesuffix(','), 16)
    i += 2
    y = int(parts[i].removesuffix(','), 16)
    i += 2
    p = int(parts[i].removesuffix(','), 16)
    i += 2
    sp = int(parts[i].removesuffix(','), 16)
    i += 2
    address = int(parts[i].removesuffix(','), 16)
    i += 2
    clock = int(parts[i])

    return Instruction(
        name=name,
        address=address,
        a=a,
        x=x,
        y=y,
        p=p,
        sp=sp,
        clock=clock,
    )


if __name__ == '__main__':
    main()
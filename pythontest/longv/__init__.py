a = 100
b = 200
print(id(a))
print(id(b))


def f(a, b):
    print(id(a))
    print(id(b))
    a = 1
    b = 2
    print(a, b)


f(a, b)

print(a, b)

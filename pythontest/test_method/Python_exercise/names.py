# from .name_function import get_formatted_name
from Python_exercise import name_function

print("Enter 'q' at any time to quit.")
while True:
    first = input("\nPlease give me a first name: ")
    if first == 'q':
        break
    # middle = input("\nPlease give me a first name: ")
    # if first == 'q':
    #     break
    last = input("Please give me a last name: ")
    if last == 'q':
        break
    formatted_name = name_function.get_formatted_name(first, last)
    print("\tNeatly formatted name: " + formatted_name + '.')

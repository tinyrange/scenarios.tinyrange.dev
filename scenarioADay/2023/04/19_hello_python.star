vm1 = main.add_vm(
    "lesson_vm",
    ram = 256,
    init_script = """
# Create student user
adduser student --disabled-password --shell /bin/ash

# Install Python
apk add python3

# Install Nano (text editor)
apk add nano
""",
    additional_files = {
        "/etc/motd": file_contents("""
Welcome to Hello, Python.

Type 'python3` to start Python.

Type 'python3 hello.py' to run a simple script.

Type 'nano hello.py' to edit the script.

"""),
        "/home/student/hello.py": file_contents("""
print("Hello, World")
"""),
    },
)

main.add_service(lesson(
    title = "Hello, Python",
    text = """
This is a simple lesson showing a VM with Python installed.
""",
    interaction = interaction_ssh(
        vm1,
        username = "student",
    ),
    open_browser = True,
))
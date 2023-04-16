"""
Testing Lesson for TinyRange
Joshua D. Scarsbrook - The University of Queensland
"""

vm1 = main.add_vm(
    "lesson_vm",
    init_script = """
# Create student user
adduser student --disabled-password --shell /bin/ash
""",
    additional_files = {
        "/etc/motd": file_contents("""
Hello World. This is a simple test lesson.

"""),
    },
)

main.add_service(lesson(
    title = "Hello, World",
    text = """
This is a testing lesson.
""",
    interaction = interaction_ssh(
        vm1,
        username = "student",
    ),
    open_browser = True,
))
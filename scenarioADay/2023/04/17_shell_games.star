vm1 = main.add_vm(
    "lesson_vm",
    init_script = """
# Create student user
adduser student --disabled-password --shell /bin/ash

echo "Hello" > /home/student/README.txt

mkdir /challenge/
mkdir /challenge/park
mkdir /challenge/park/tree
mkdir /challenge/house/bedroom
mkdir /challenge/house/bedroom/bed
mkdir /challenge/house/bedroom/closet
echo "flag{congratulations}" > /challenge/house/bedroom/closet/flag.txt
mkdir /challenge/house/lounge
mkdir /challenge/house/lounge/tv
mkdir /challenge/house/bathroom
mkdir /challenge/house/kitchen
mkdir /challenge/house/kitchen/panty
""",
    additional_files = {
        "/etc/motd": file_contents("""
Welcome to Hide and Seak.

"""),
    },
)

main.add_service(lesson(
    title = "Hide and Seak",
    text = """
Welcome to Hide and Seak.

Your goal is to find the flag somewhere on this filesystem.

The format of the flag is `flag{<some text here>}`.

**Hint:** The flag is somewhere underneath `/challenge`.

**Hint:** Use these commands to find the flag.

- `pwd`: Print the current working directory.
- `cd`: Change to a different directory. Try `cd /challenge`.
- `cat`: Read a file. Try `cat README.txt`.
""",
    interaction = interaction_ssh(
        vm1,
        username = "student",
    ),
    open_browser = True,
))
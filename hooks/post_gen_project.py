import os
import re
import shutil
import stat

import textwrap
from pathlib import Path
from sys import exit
from typing import Callable, List, Optional, Any

# Contains path to the project directory - i.e. {{cookiecutter.project_name.strip()}}
CUR_DIR = Path.cwd().absolute()


def remove_precommit():
    """
    If pre-commit is not to be used, removes pre-commit configuration file, and the Github workflow
    for the same
    """

    use_precommit: bool = bool("{{ cookiecutter.use_precommit }}".lower() == "y")
    if use_precommit:
        return

    config_file: str = os.path.join(CUR_DIR, ".pre-commit-config.yaml")
    workflow: str = os.path.join(CUR_DIR, ".github", "workflows", "pre-commit.yml")

    if os.path.exists(config_file):
        os.remove(config_file)

    if os.path.exists(workflow):
        os.remove(workflow)


def check_email_provided():
    """
    Checks if the user has provided an email while generating the template

    Removes the `SECURITY.md` and `CODE_OF_CONDUCT.md` file in case the email field is empty -- if
    the email ID is empty, people cannot get in touch with project stakeholders to report, making
    the file(s) be pointless.
    """

    should_remove: bool = bool("{{ cookiecutter.contact_email }}" == "")
    if not should_remove:
        return

    # List of files to be removed
    files: List[str] = ["SECURITY.md", "CODE_OF_CONDUCT.md"]
    for f in files:
        deletable_file = os.path.join(CUR_DIR, f)
        if not os.path.exists(path=deletable_file):
            continue

        os.remove(deletable_file)


def print_final_instructions():
    """
    Simply prints final instructions for users to follow once they generate a project
    using this template!
    """

    message = """
    ====================================================================================

    Your project `{{ cookiecutter.project_name.strip() }}` is ready!

    It is suggested to go through the generated `README` for detailed instructions on
    how to use the tools that come with `go-template`! The following is a *brief*
    overview of steps to push code to remote.

    1) Move to project directory, and initialize a git repository:
        $ cd {{ cookiecutter.project_name.strip() }} && git init

    2) Configure GOPRIVATE if using private Go modules
        $ go env -w GOPRIVATE="github.com/{{ cookiecutter.org }}/*"

    3) Run proto generation:
        $ make proto

    4) Run linters and codestyle checkers
        $ make lint

    5) Run the complete test-suite:
        $ make test

    6) Upload initial code to git:
        $ git add -a
        $ git commit -m "Initial commit!"
        $ git remote add origin https://{{ cookiecutter.go_module_path.strip('/') }}.git
        $ git push -u origin --all

    If you have any suggestions, or encounter any issue, please raise an issue at your template repository.
    """

    print(textwrap.dedent(message))


def make_scripts_executable():
    """
    Ensures scripts are executable after cookiecutter generation.
    """

    scripts_dir = CUR_DIR / "scripts"
    if not scripts_dir.exists():
        return

    for script in scripts_dir.rglob("*.sh"):
        if not script.is_file():
            continue
        mode = script.stat().st_mode
        script.chmod(mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)


def remove_empty_workflow():
    """
    Removes the `.github/workflows` directory if it is empty
    """

    workflow_dir = os.path.join(CUR_DIR, ".github", "workflows")
    if not os.path.exists(workflow_dir):
        return
    # iterate over the files in the directory
    files = os.listdir(workflow_dir)
    for f in files:
        # check if the file is empty
        path = os.path.join(workflow_dir, f)
        if os.path.getsize(path) <= 4:
            os.remove(path)


runners: Callable[[Optional[Any]], None] = [
    remove_precommit,
    check_email_provided,
    make_scripts_executable,
    print_final_instructions,
    remove_empty_workflow,
]

for runner in runners:
    try:
        runner()
    except ValueError as e:
        print(e)
        exit(-10)

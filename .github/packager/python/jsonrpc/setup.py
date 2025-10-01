from pathlib import Path

from setuptools import find_packages, setup

setup(
    name="pactus-jsonrpc",
    version="{{ VERSION }}",
    author="Pactus Development Team",
    author_email="info@pactus.org",
    url="https://pactus.org",
    description="Python client for interacting with the Pactus blockchain via JSON-RPC",
    long_description=Path("README.md").read_text(encoding="utf-8"),
    long_description_content_type="text/markdown",
    url="https://github.com/pactus-project/pactus",
    packages=find_packages(),
    license="MIT",
    install_requires=[
        "jsonrpc2-pyclient>=5.2.0",
        "py-undefined>=0.1.5",
        "pydantic>=2.5.3"
    ],
    keywords=["pactus", "blockchain", "json-rpc"],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Build Tools",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)

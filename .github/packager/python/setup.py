from setuptools import setup, find_packages

setup(
    name="pactus-grpc",
    version="{{ VERSION }}",
    author="Pactus Development Team",
    author_email="info@pactus.org",
    description="gRPC client bindings for the Pactus blockchain",
    long_description=open("README.md", encoding="utf-8").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/pactus-project/pactus",
    license="MIT",
    packages=find_packages(),
    keywords=["pactus", "blockchain", "grpc"],
    install_requires=[
        "grpcio",
        "protobuf",
    ],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Build Tools",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)

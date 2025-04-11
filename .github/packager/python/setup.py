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
    packages=find_packages(),
    install_requires=[
        "grpcio",
        "protobuf",
    ],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)

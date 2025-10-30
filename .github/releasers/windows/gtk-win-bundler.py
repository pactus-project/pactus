#!/usr/bin/env python3
"""
GTK Bundle Helper for Windows
Automatically detects and bundles all GTK dependencies for Windows distribution.
"""

import os
import sys
import shutil
import subprocess
import json
from pathlib import Path
from typing import List, Set, Dict, Optional

class GTKBundler:
    def __init__(self, mingw_prefix: str, target_dir: str):
        self.mingw_prefix = Path(mingw_prefix)
        self.target_dir = Path(target_dir)

    def run_command(self, cmd: List[str]) -> str:
        """Run a command and return its output."""
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, check=True)
            return result.stdout.strip()
        except subprocess.CalledProcessError as e:
            print(f"Error running command {' '.join(cmd)}: {e}")
            sys.exit(1)

    def get_dependencies(self, exe_path: Path, dependencies: List[Path]) -> None:
        """Get all DLL dependencies for an executable using ldd."""

        # Use ldd to get dependencies
        ldd_output = self.run_command(['ldd', str(exe_path)])

        for line in ldd_output.split('\n'):
            if '/mingw64' in line and '.dll' in line:
                # Extract the DLL path
                parts = line.split()
                if len(parts) >= 3:
                    dll_path = parts[2].replace('/mingw64', str(self.mingw_prefix))
                    dll_path = Path(dll_path)
                    if dll_path not in dependencies:
                        dependencies.append(dll_path)

    def copy_file(self, src: Path, dst: Path) -> None:
        """Copy a file if it hasn't been copied already."""
        if dst.exists():
            print(f"  Already copied: {src.name}")
            return

        if not src.exists():
            print(f"ERROR: Required file not found: {src}")
            sys.exit(1)

        try:
            # Create parent directories if needed
            dst.parent.mkdir(parents=True, exist_ok=True)

            # Copy the file
            shutil.copy2(src, dst)
            print(f"  Copied: {src.name}")
        except Exception as e:
            print(f"ERROR: Failed to copy file {src} to {dst}: {e}")
            sys.exit(1)

    def copy_dir(self, src_path: str, dst_path: str) -> None:
        """Copy a directory if it exists."""
        if not Path(src_path).exists():
            print(f"ERROR: Required directory not found: {src_path}")
            sys.exit(1)

        try:
            Path(dst_path).parent.mkdir(parents=True, exist_ok=True)
            shutil.copytree(src_path, dst_path, dirs_exist_ok=True)
            print(f"  Copied directory: {Path(src_path).name}")
        except Exception as e:
            print(f"ERROR: Failed to copy directory {src_path} to {dst_path}: {e}")
            sys.exit(1)

    def copy_gtk_dependencies(self, exe_path: Path, target_exe_dir: Path) -> None:
        """Copy all dependencies for the GTK executable."""
        print(f"Analyzing dependencies for: {exe_path.name}")

        # Get direct dependencies
        dependencies = []
        self.get_dependencies(exe_path, dependencies)

        for dep in dependencies:
            self.get_dependencies(dep, dependencies)

        # Add dependencies for all pixbuf loader DLLs
        pixbuf_dir = Path(f"{self.mingw_prefix}/lib/gdk-pixbuf-2.0/2.10.0/loaders")
        for loader_file in pixbuf_dir.glob("*.dll"):
            print(f"    Scanning pixbuf loader: {loader_file.name}")
            self.get_dependencies(loader_file, dependencies)

        for dep in dependencies:
            self.copy_file(dep, target_exe_dir / dep.name)

    # Copy GTK resources
    # Based on this tutorial: https://www.gtk.org/docs/installations/windows#building-and-distributing-your-application
    def copy_gtk_resources(self) -> None:
        """Copy GTK resources (themes, icons, schemas, etc.)."""
        print("Copying GTK resources...")

        # Copy GdkPixbuf loaders
        print("  Copying GdkPixbuf loaders...")
        self.copy_dir(
            f"{self.mingw_prefix}/lib/gdk-pixbuf-2.0",
            f"{self.target_dir}/lib/gdk-pixbuf-2.0"
        )

        # Copy icon themes
        print("  Copying icon themes...")
        self.copy_dir(
            f"{self.mingw_prefix}/share/icons/Adwaita",
            f"{self.target_dir}/share/icons/Adwaita"
        )
        self.copy_dir(
            f"{self.mingw_prefix}/share/icons/hicolor",
            f"{self.target_dir}/share/icons/hicolor"
        )

        # Copy GTK themes
        print("  Copying GTK themes...")
        self.copy_dir(
            f"{self.mingw_prefix}/share/gtk-3.0",
            f"{self.target_dir}/share/themes/Windows10/gtk-3.0"
        )

        # Copy GSettings schemas
        print("  Copying GSettings schemas...")
        self.copy_dir(
            f"{self.mingw_prefix}/share/glib-2.0/schemas",
            f"{self.target_dir}/share/glib-2.0/schemas"
        )

        # Create settings.ini
        settings_file = f"{self.target_dir}/share/gtk-3.0/settings.ini"
        Path(settings_file).parent.mkdir(parents=True, exist_ok=True)
        with open(settings_file, 'w') as f:
            f.write("[Settings]\n")
            f.write("gtk-theme-name=Adwaita\n")
            f.write("gtk-icon-theme-name=Adwaita\n")
            f.write("gtk-font-name=Segoe UI 9\n")
            f.write("gtk-application-prefer-dark-theme=true\n")
        print("  Created settings.ini")

    def copy_gtk_executables(self, target_exe_dir: Path) -> None:
        """Copy helper executables that might be needed."""
        print("Copying helper executables...")

        helper_exes = [
            "gdbus.exe",
            "gspawn-win64-helper.exe",
            "gspawn-win64-helper-console.exe"
        ]

        for exe in helper_exes:
            exe_src = self.mingw_prefix / "bin" / exe
            if exe_src.exists():
                self.copy_file(exe_src, target_exe_dir / exe)


    def bundle_application(self, exe_path: Path) -> None:
        """Main method to bundle the entire application."""
        print(f"Bundling GTK application: {exe_path}")

        # Create target directory structure
        target_exe_dir = self.target_dir
        target_exe_dir.mkdir(parents=True, exist_ok=True)

        # Copy all dependencies
        self.copy_gtk_dependencies(exe_path, target_exe_dir)

        # Copy helper executables
        self.copy_gtk_executables(target_exe_dir)

        # Copy GTK resources
        self.copy_gtk_resources()

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 gtk-win-bundler.py <exe_path> <target_dir>")
        sys.exit(1)

    # Check for MINGW_PREFIX environment variable
    mingw_prefix = os.environ.get('MINGW_PREFIX')
    if not mingw_prefix:
        print("Error: MINGW_PREFIX environment variable is not set")
        sys.exit(1)

    exe_path = Path(sys.argv[1])
    target_dir = Path(sys.argv[2])

    bundler = GTKBundler(mingw_prefix, target_dir)
    bundler.bundle_application(exe_path)

if __name__ == "__main__":
    main()

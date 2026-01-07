# Generating the Windows resource (.syso)
1) Install the tool (once):
   `go install github.com/akavel/rsrc@latest`

2) Generate the .syso in this folder:
   `rsrc -manifest ./pactus-gui.manifest -ico ./pactus.ico`

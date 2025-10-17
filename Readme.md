# Pennsieve Processor Pre-Packages v2 (Folder Rebuild)

A Go-based preprocessor service that downloads and organizes files from Pennsieve datasets while preserving their original folder structure. This project is based on the original [processor-pre-packages-v2](https://github.com/Pennsieve/processor-pre-packages-v2).

## Purpose

This service prepares files from a Pennsieve workspace dataset for downstream processing by analytical tools. Unlike the standard processor-pre-packages-v2 which flattens all files into a single directory, this version **maintains the original folder hierarchy** of your files.

## What's Different?

| Feature | Standard Version | Folder Rebuild Version (This Repo) |
|---------|------------------|-------------------------------------|
| File Organization | All files flattened to single directory | **Original folder structure preserved** |
| Use Case | Simple processors that don't need structure | Processors requiring specific directory access |
| Output | `OUTPUT_DIR/file1.txt`, `OUTPUT_DIR/file2.csv` | `OUTPUT_DIR/folder1/file1.txt`, `OUTPUT_DIR/folder2/subfolder/file2.csv` |


## Need Help?

- **Original Project**: [processor-pre-packages-v2](https://github.com/Pennsieve/processor-pre-packages-v2)
- **Pennsieve Documentation**: [docs.pennsieve.io](https://docs.pennsieve.io)





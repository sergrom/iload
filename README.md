[![build](https://github.com/sergrom/iload/workflows/build/badge.svg)](https://github.com/sergrom/iload/actions/workflows/build.yml)

# iload
Multithread file downloader

## Install
1. First, install golang:
https://go.dev/doc/install

2. Then type in console:
```bash
$ go install github.com/sergrom/iload
```

## Usage
To see parameters and options type:
```bash
$ iload -h
```

## Example
If you want to download files with urls listed in <code>/path/to/input/file.txt</code> to directory <code>/directory/to/output/</code> using <code>5 threads</code> and with <code>verbose</code> mode type:
```bash
$ iload -f /path/to/input/file.txt -d /directory/to/output/ -t 5 -v
````
or
```bash
iload --input-file="/path/to/input/file.txt" --output-dir="/directory/to/output/" --threads-num=5 --verbose
```

## Parameters and options
<code>-h</code>, <code>--help</code>                help for iload<br>
<code>-f</code>, <code>--input-file</code> string   Input file with urls you want to download. Each url must be in separate line<br>
<code>-d</code>, <code>--output-dir</code> string   Output directory to which you want to save downloaded files<br>
<code>-t</code>, <code>--threads-num</code> int     Number of threads (default 5)<br>
<code>-v</code>, <code>--verbose</code>             Verbose

## Input file
The input file must be specified by parameter <code>-f</code> and must contain a list of urls to files.
Each url should be in separate line like this:
<pre>
http://site.com/path/to/file1.jpg
http://site.com/path/to/file2.jpg
http://site.com/path/to/file3.jpg
</pre>
The program will try to download and save files with names: "file1.jpg", "file2.jpg", "file3.jpg" to output directory.


If you want to save files with another names, you can specify them like this:
<pre>
http://site.com/path/to/file1.jpg|first_file.jpg
http://site.com/path/to/file2.jpg|second_file.jpg
http://site.com/path/to/file3.jpg|third_file.jpg
</pre>
The program will try to download and save files with names: "first_file.jpg", "second_file.jpg", "third_file.jpg" to output directory.

If while saving some file with name <code>file.jpg</code> it happens that file with same name is already exists,
the program wil try to save it with name <code>file_1.jpg</code>.
If such file also exists, it will try to save the file with name <code>file_2.jpg</code>, and so on unil <code>file_99.jpg</code>, then program will start to rewrite files in output directory. So if you want to prevent rewriting downloaded files, please specify file names explicitly.
## Output directory
The output directory must be specified by parameter <code>-d</code>.

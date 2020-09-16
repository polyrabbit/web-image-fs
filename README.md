# imagefs - Mount web images as local files

## Example

View my blog in your explorer/finder:

```bash
$ imagefs -v /tmp/imagefs https://blog.betacat.io/post/2020/08/how-to-mount-etcd-as-a-filesystem/
$ ls -alh
total 648K
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 ./
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 ↩︎/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 1/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 2/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 3/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 About/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Archives/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 CC BY-NC-ND 4.0/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 dentry cache/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 EntryTimeout/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 etcdfs/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 FUSE/
-rw-r--r-- 0 poly poly 166K Sep 11 18:01 FUSE Stack.png
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 FUSE 文件系统/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 go-fuse/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 go-fuse Inode structure/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Go-fuse 库/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Home/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Hugo/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 inode/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 inode编号/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Jane/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 Kernel FUSE message format/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 OpCode enum/
-rw-r--r-- 0 poly poly 217K Sep 11 18:01 open kubernetes etcd in vscode.png
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 sshfs/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 TL;DR/
-rw-r--r-- 0 poly poly 257K Sep 11 18:01 VFS Read Operation.png
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 举个栗子/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 喵叔/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 喵叔没话说/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 工作原理/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 总结/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 文件系统相关的代码/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 沪ICP备17033881号-1/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 编写可测试 Go 代码的一种模式?            下一篇/
drwxr-xr-x 0 poly poly    0 Sep 11 18:01 背景/
```

## Usage

```bash
$ imagefs
Mount web images to local file system - find help/update at https://github.com/polyrabbit/web-image-fs

Usage:
  imagefs [mount-point] [url] [flags]

Flags:
      --http-timeout duration   http request timeout (default 10s)
      --enable-pprof            enable runtime profiling data via HTTP server. Address is at "http://localhost:9327/debug/pprof"
  -v, --verbose                 verbose output
      --mount-options strings   options are passed as -o string to fusermount (default [nonempty])
  -h, --help                    help for imagefs
```

## Limitations

This tool uses a simple html parser called [goquery](https://github.com/PuerkitoBio/goquery) that does not evaluate javascript, so it cannot handle dynamically generated images.

## TODO

- [x] ~~Prefetch image stats using a `HEAD` request~~
- [ ] Detect file extension using `content-type` header
- [ ] Use a cache library for `DomNode` objects
- [x] ~~`DomNode` should contain an absolute self link url~~

## License

The MIT License (MIT) - see [LICENSE.md](https://github.com/polyrabbit/web-image-fs/blob/master/LICENSE) for more details

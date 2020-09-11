# imagefs - Mount web images as local files

## Example

View my blog in your explorer/finder:

```bash
$ imagefs -v /tmp/imagefs https://blog.betacat.io/post/2020/08/how-to-mount-etcd-as-a-filesystem/
$ ls -alh
total 0
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 ./
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 ../
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 1/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 2/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 About/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 Archives/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 EntryTimeout/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 FUSE/
-rw-r--r--  0 bytedance  staff   166K Sep 11 08:23 FUSE Stack.png
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 Home/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 Hugo/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 Kernel FUSE message format/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 TL;DR/
-rw-r--r--  0 bytedance  staff   256K Sep 11 08:23 VFS Read Operation.png
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 go-fuse/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 go-fuse Inode structure/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 inode/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 inode编号/
-rw-r--r--  0 bytedance  staff   216K Sep 11 08:23 open kubernetes etcd in vscode.png
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 沪ICP备17033881号-1/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 喵叔/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 总结/
drwxr-xr-x  0 bytedance  staff     0B Sep 11 08:23 文件系统相关的代码/
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
- [ ] `DomNode` should contain an absolute self link url

## License

The MIT License (MIT) - see [LICENSE.md](https://github.com/polyrabbit/web-image-fs/blob/master/LICENSE) for more details

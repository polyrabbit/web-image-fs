# imagefs - Mount web images as local files

## Example

View my blog in your explorer/finder:

```bash
$ imagefs -v /tmp/imagefs https://blog.betacat.io/post/2020/08/how-to-mount-etcd-as-a-filesystem/
$ ls -alh
total 0
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 ./
drwxrwxrwt  17 root       wheel   544B Sep 10 23:23 ../
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 3/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 About/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 Archives/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 Go-fuse 库/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 Home/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 Jane/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 dentry cache/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 etcdfs/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 go-fuse Inode structure/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 inode/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 inode编号/
-rw-r--r--   0 bytedance  staff     0B Sep 10 23:41 open kubernetes etcd in vscode.png
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 sshfs/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 沪ICP备17033881号-1/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 ↩︎/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 喵叔/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 总结/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 背景/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 举个栗子/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 喵叔没话说/
drwxr-xr-x   0 bytedance  staff     0B Sep 10 23:41 编写可测试 Go 代码的一种模式?            下一篇/
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

- [ ] Prefetch image stats using a `HEAD` request
- [ ] Detect file extension using `content-type` header
- [ ] Use a cache library for `DomNode` objects
- [ ] `DomNode` should contain an absolute self link url

## License

The MIT License (MIT) - see [LICENSE.md](https://github.com/polyrabbit/web-image-fs/blob/master/LICENSE) for more details

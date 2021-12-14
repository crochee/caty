#!/bin/bash
set -e

#rust-analyzer加载时间过长
#
#开发环境:
#vscode+rust-analyzer
#问题：
#vscode一直卡在fetching metadata阶段。
#方法:
#运行cargo metadata,发现
#
#Blocking waiting for file lock on package cache
#运行
#
#rm -rf ~/.cargo/.package-cache
#删除cargo的缓存，而后再运行cargo metadata后发现没有blocking的提示，重启vscdoe，加载rust-analyzer成功。

#curl -L https://github.com/rust-analyzer/rust-analyzer/releases/latest/download/rust-analyzer-x86_64-unknown-linux-gnu.gz | gunzip -c - > ~/.go/bin/rust-analyzer
# chmod + ~/.cargo/bin/rust-analyzer
# { "rust-analyzer.server.path": "~/.cargo/bin/rust-analyzer" }
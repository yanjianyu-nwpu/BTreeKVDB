# BTreeKVDB
简单的b+kv数据库

## 事务
### 支持同时一个写事务和多个读事务
### 多个读事务采用MVCC实现

## MVCC的实现
### 在读事务未提交时候会保留该版本需要的页在内存中。

## 数据存储结构
### 全部采用sync模式，使用B+树存储。commit时候将页一定落盘。

## 页面结构
### 页面分为根节点、branch节点和叶子节点。数据放在叶子节点中。

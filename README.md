# agones_sync_interval_test

[Agones](https://github.com/googleforgames/agones) の [CustomFasSyncInterval](https://github.com/googleforgames/agones/issues/1955) をローカル環境で試しに動かしてみただけのもの。

* CustomFasSyncInterval を有効にして Agones をインストールして
* Sync Interval を変えて FleetAutoscaler を立てて
* allocate を何回か実施した

だけ。
テストで検証しているように Sync Interval の変更が有効になっていることが確認できる。

このリポジトリに置いたテストのような fleet が枯渇した後の挙動を確かめるテストを書きたいときに Sync Interval をいじれるとテストがすぐに終わってハッピーになれそう。
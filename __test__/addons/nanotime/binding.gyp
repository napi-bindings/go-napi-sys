{
  "targets": [
    {
      "target_name": "addon",
      "sources": [ "addon.cc" ],
      "include_dirs": [
        "include"
      ],
      "libraries": [
        "<(module_root_dir)/libgoaddon.a"
      ]
    }
  ]
}
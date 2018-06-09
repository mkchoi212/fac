# ⚙️  Configuration

## .fac.yml

Configure fac by adding a `.fac.yml` in your `$HOME` directory. With the file in place, the following paramters can be configued:

```yml
select_local: a
select_incoming: b
toggle_view: c
show_up: d
show_down: e 
scroll_up: f
scroll_down: g
edit: h
next: i
previous: j
quit: k
help: l
```

## 📖  Rules

When parsing `.fac.yml`, fac enforces *three rules.*

### 1. Invalid key-binding keys
```yml
foobar: f
```

> Warning: Invalid key: "fobar" will be ignored

### 2. Multi-character key-mappings
```yml
select_local: local
```

> Warning: Illegal multi-character mapping: "local" will be interpreted as 'l'  

### 3. Duplicate key-mappings
```yml
show_up: u
scroll_up: u
```

> Fatal: Duplicate key-mapping: "scroll_up, show_up" are all represented by 'u'

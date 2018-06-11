# HOW TO: Custom Configuration

## .fac.yml

You can configure fac by creating a `.fac.yml` in your `$HOME` directory. With the file in place, the following paramters can be tweaked:

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

# Set to `true` to skip having to press enter when entering commands
cont_eval: true
```

## 👨‍⚖️👩‍⚖️ Rules 

When parsing `.fac.yml`, fac enforces *three rules.*

### 1. Invalid key-binding keys

```yml
# WRONG
# Warning: Invalid key: "foobar" will be ignored
foobar: f
```

### 2. Multi-character key-mappings

```yml
# WRONG
# Warning: Illegal multi-character mapping: "local" will be interpreted as 'l'  
select_local: local

# CORRECT
select_local: l
```

### 3. Duplicate key-mappings

```yml
# WRONG
# Fatal: Duplicate key-mapping: "scroll_up, show_up" are all represented by 'u'
show_up: u
scroll_up: u

# CORRECT
show_up: u
scroll_up: k
```

### 4. Mapping keys

```yml
# WRONG
# Warning: yaml: mapping keys are not allowed in this context
help: ?

# CORRECT
help: "?"
```

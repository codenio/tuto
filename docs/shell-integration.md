# Shell prompt integration

`tuto session prompt` prints a compact token like `(bash-scripting: 2/10)` when a session
is active, and **nothing** when there is no session. Drop one line into your shell config to
make it appear in your prompt automatically.

---

## Zsh

```zsh
# ~/.zshrc
RPROMPT='$(tuto session prompt 2>/dev/null)'
```

For left-side prompt instead:

```zsh
PROMPT='$(tuto session prompt 2>/dev/null) '$PROMPT
```

---

## Bash

```bash
# ~/.bashrc
PS1='$(tuto session prompt 2>/dev/null) '$PS1
```

---

## Fish

```fish
# ~/.config/fish/functions/fish_right_prompt.fish
function fish_right_prompt
    tuto session prompt 2>/dev/null
end
```

---

## Starship

Add a custom command module to `~/.config/starship.toml`:

```toml
[custom.tuto]
command = "tuto session prompt"
when = "tuto session prompt"
format = "[$output]($style) "
style = "bold cyan"
```

---

## Example output

```
➜  myproject  (bash-scripting: 2/10)
```

The token disappears automatically when you run `tuto session end` or complete a module.

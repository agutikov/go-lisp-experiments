package cmdlex

import (
)


type CmdLineArgs struct {
	Flags map[string]uint
	Options map[string][]string
	Positional []string

	Invalids []string

	Next *CmdLineArgs // TODO levels
}

func newCmdLineArgs() CmdLineArgs {
	return CmdLineArgs{
		Flags: map[string]uint{},
		Options: map[string][]string{},
	}
}


func (e *CmdLineArgs) add_invalid(s string) {
	e.Invalids = append(e.Invalids, s)
}

func (e *CmdLineArgs) add_positional(s string) {
	e.Positional = append(e.Positional, s)
}

func (e *CmdLineArgs) add_flag(s string) {
	count, ok := e.Flags[s]
	if ok {
		e.Flags[s] = count + 1
	} else {
		e.Flags[s] = 1
	}
}

func (e *CmdLineArgs) add_option(key string, value string) {
	lst, ok := e.Options[key]
	if ok {
		e.Options[key] = append(lst, value)
	} else {
		e.Options[key] = []string{value}
	}
}

type arg_desc struct {
	// A little bit redundant but handy
	valid bool
	dash int
	key string
	eq bool
	value string
}

func (d *arg_desc)is_key() bool {
	return d.valid && d.dash >= 1 && d.dash <= 2 && len(d.key) > 0 && !d.eq
}

func (d *arg_desc)is_value() bool {
	return d.valid && len(d.value) > 0 && len(d.key) == 0
}

func (d *arg_desc)is_kv() bool {
	return d.valid && len(d.value) > 0 && len(d.key) > 0 && d.eq
}

func new_arg_desc() arg_desc {
	return arg_desc{false, 0, "", false, ""}
}

func parse_arg(s string) arg_desc {
	d := new_arg_desc()

	key_first_index := 0

	for i, c := range s {
		if c == '-' && i == key_first_index {
			// leading dashes
			key_first_index++
			d.dash++
			continue
		}
		if key_first_index > 2 {
			// too many leading dashes
			return d
		}
		if c == '=' {
			if i <= key_first_index {
				// leading =
				return d
			} else if d.dash > 0 {
				// [-[-]]key=value
				//TODO: single key single dash
				d.eq = true
				d.key = s[0:i-1]
				d.value = s[i+1:]
				d.valid = len(d.value) > 0
				return d
			}
		}
	}

	if d.dash > 0 && key_first_index <= len(s)-1 {
		// -[-]key
		d.key = s[key_first_index:]
	} else {
		// value
		d.value = s
	}
	d.valid = d.valid || len(d.key) > 0 || len(d.value) > 0

	return d
}

func ParseCmdLineArgs(args []string, levels int) CmdLineArgs {
	r := newCmdLineArgs()

	if len(args) == 0 {
		return r
	}

	last := &r

	after_double_dash := false
	prev_arg_desc := new_arg_desc()

	r.add_positional(args[0])

	for i, arg := range args[1:] {
		if after_double_dash {
			last.add_positional(arg)
			continue
		}
		if arg == "--" {
			after_double_dash = true

			if prev_arg_desc.is_key() {
				r.add_flag(prev_arg_desc.key)
			}

			continue
		}

		curr := parse_arg(arg)
		//fmt.Printf("%+v\n", curr)

		if curr.is_key() && i == len(args)-2 {
			// trailing flag - bool
			r.add_flag(curr.key)
		}

		if prev_arg_desc.is_key() {
			if curr.is_key() || curr.is_kv() {
				// previous is flag
				r.add_flag(prev_arg_desc.key)
			} else {
				r.add_option(prev_arg_desc.key, arg)
			}
		} else {
			if !curr.valid {
				r.add_invalid(arg)
			} else if curr.is_value() {
				r.add_positional(arg)
			} else if curr.is_kv() {
				r.add_option(curr.key, curr.value)
			} else {
				// current is a key
			}
		}

		prev_arg_desc = curr
	}

	return r
}





<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="left" height="96" width="96"> Extension
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a>::<span>Common</span> toolkit.</h6>


## 🌱 Overview

The **common/extension** module provides parsing and conversion utilities for core data types.  
Each subpackage focuses on a specific type (boolean, date, decimal, float, integer, number, string),  
offering safe parsing functions with consistent error handling.

---

## 📦 Subpackages

| Package                | Description                                                                                         |
|------------------------|-----------------------------------------------------------------------------------------------------|
| [`boolean`](./boolean) | Parses arbitrary inputs into booleans (`true/false`, `1/0`, `yes/no`, `on/off`, `y/n`, `t/f`).      |
| [`date`](./date)       | Date parsing utilities with cleaning, normalization, and leap year validation.                      |
| [`decimal`](./decimal) | Parses values into `float64` with a specified precision (`3.14159` → `3.14` if precision=2).        |
| [`float`](./float)     | Parses values into `float64` with no restrictions. Supports ints, floats, strings, bools, pointers. |
| [`integer`](./integer) | Parses values into `int`, always truncating toward zero (`3.99` → `3`, `-1.9` → `-1`).              |
| [`number`](./number)   | Numeric parser layer that supports float or int with optional rounding flag.                        |
| [`string`](./string)   | String parsing and normalization utilities. Converts values to strings safely.                      |

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project

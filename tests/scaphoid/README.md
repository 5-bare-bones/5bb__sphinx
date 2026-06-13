
# Summary report on 2026-06-05

  ┌──────────────────────────┬───────┬──────────────────────────────────────────────────────────────────────────────────┐
  │           File           │ Tests │                                  What's covered                                  │
  ├──────────────────────────┼───────┼──────────────────────────────────────────────────────────────────────────────────┤
  │ test_streams.py          │ 71    │ OutputStream fields, _fail format, all 6 assertions + accessors + dunders;       │
  │                          │       │ StreamSearch fields, _fail format, all 3 methods                                 │
  ├──────────────────────────┼───────┼──────────────────────────────────────────────────────────────────────────────────┤
  │ test_result.py           │ 47    │ from_completed all 7 fields (incl. cwd default + copy semantics), ok, output,    │
  │                          │       │ _context all 6 parts, _fail, all 3 assertions + chaining across levels           │
  ├──────────────────────────┼───────┼──────────────────────────────────────────────────────────────────────────────────┤
  │                          │       │ RunConfig all 7 defaults + individual fields + validation + replace;             │
  │ test_binary_execution.py │ 60    │ _resolve_env all 3 branches; run() command forms, output capture, exit codes,    │
  │                          │       │ metadata, check, env injection; Shell init, run, set_env                         │
  ├──────────────────────────┼───────┼──────────────────────────────────────────────────────────────────────────────────┤
  │                          │       │ before_scenario all 5 context fields + tmpdir on disk + prefix + uniqueness;     │
  │ test_environment.py      │ 16    │ after_scenario removes dir + with files + with subdirs + missing attribute +     │
  │                          │       │ already deleted + path is a file                                                 │
  ├──────────────────────────┼───────┼──────────────────────────────────────────────────────────────────────────────────┤
  │ test_imports.py          │ 28    │ Public API surface, __all__ completeness, module importability, re-export        │
  │                          │       │ identity                                                                         │
  └──────────────────────────┴───────┴──────────────────────────────────────────────────────────────────────────────────┘

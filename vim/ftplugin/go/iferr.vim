scriptencoding utf-8

function! s:IfErr()
  let bpos = wordcount()['cursor_bytes']
  let out = systemlist('iferr -pos ' . bpos, bufnr('%'))
  if len(out) == 1
    return
  endif
  let pos = getcurpos()
  call append(pos[1], out)
  silent normal! j=2j
  call setpos('.', pos)
  silent normal! 4j
endfunction

command! -buffer -nargs=0 IfErr call s:IfErr()

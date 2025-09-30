#ifndef LUAWRAP_H
#define LUAWRAP_H

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

// Wrapper functions for Lua macros that cgo can't handle directly

static inline void luawrap_pop(lua_State *L, int n) {
    lua_pop(L, n);
}

static inline int luawrap_pcall(lua_State *L, int nargs, int nresults, int errfunc) {
    return lua_pcall(L, nargs, nresults, errfunc);
}

static inline int luawrap_isfunction(lua_State *L, int idx) {
    return lua_isfunction(L, idx);
}

static inline lua_Number luawrap_tonumber(lua_State *L, int idx) {
    return lua_tonumber(L, idx);
}

static inline const char* luawrap_tostring(lua_State *L, int idx) {
    return lua_tostring(L, idx);
}

#endif // LUAWRAP_H

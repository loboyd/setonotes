To make the security of setonotes easier to reason about (and easier to audit),
the code in this package should be the only code allowed to handle encrypted
keys or passwords (the exceptions are the signin and and signup handlers which
must receive the plaintext passwords submitted with the forms). This should be
enforced as much as possible by carefully restricting exports from this package
as well as returned values. It should also be a guiding design principle for
security-related code. Presently, there are a few areas of the codebase which do
not strictly follow this standard, but it is a goal to refactor these.

#!/bin/bash

gb vendor fetch -no-recurse -revision 709fab3d192d7c62f86043caff1e7e3fb0f42bd8 github.com/rwcarlsen/goexif
gb vendor fetch -no-recurse -revision 351dc6a5bf92a5f2ae22fadeee08eb6a45aa2d93 golang.org/x/crypto/ssh
gb vendor fetch -no-recurse -revision daf2955e742cf123959884fdff4685aa79b63135 github.com/olekukonko/tablewriter
gb vendor fetch -no-recurse -revision d6bea18f789704b5f83375793155289da36a3c7f github.com/mattn/go-runewidth

gb vendor fetch -no-recurse -revision 075e191f18186a8ff2becaf64478e30f4545cdad golang.org/x/net/context             # https://go.googlesource.com/net/context                       
gb vendor fetch -no-recurse -revision 04e1573abc896e70388bd387a69753c378d46466 golang.org/x/oauth2                  # https://go.googlesource.com/oauth2                            
gb vendor fetch -no-recurse -revision 04e1573abc896e70388bd387a69753c378d46466 golang.org/x/oauth2/google           # https://go.googlesource.com/oauth2/google                     
gb vendor fetch -no-recurse -revision 3261f00d16e92932f49a39672dfd540896ed30d0 cloud.google.com/go/compute/metadata # https://code.googlesource.com/gocloud/compute/metadata        
gb vendor fetch -no-recurse -revision 3261f00d16e92932f49a39672dfd540896ed30d0 cloud.google.com/go/internal         # https://code.googlesource.com/gocloud/internal                
gb vendor fetch -no-recurse -revision 518eda9a0920a55ffe7190db96fe8ed85a62e376 google.golang.org/api/gensupport     # https://code.googlesource.com/google-api-go-client/gensupport 
gb vendor fetch -no-recurse -revision 518eda9a0920a55ffe7190db96fe8ed85a62e376 google.golang.org/api/gmail/v1       # https://code.googlesource.com/google-api-go-client/gmail/v1   
gb vendor fetch -no-recurse -revision 518eda9a0920a55ffe7190db96fe8ed85a62e376 google.golang.org/api/googleapi      # https://code.googlesource.com/google-api-go-client/googleapi  

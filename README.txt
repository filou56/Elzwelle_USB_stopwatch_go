tk.Button("Name",tk.Option("text","test")) -> tkk::button .Name -text "test"

type Option struct {
    name 	string
    value	{}
}

filou@sambesi:~$ wish
% option add *Font {Helvetica 32} widgetDefault 
% button .b -text "Hello, world!" -command exit
.b
% grid .b
% filou@sambesi:~$ 

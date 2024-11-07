import sys


class ArgOpts:
     def __init__(self, default: any, requires_val: bool, description: str):
        self.value = default
        self.description = description
        self.requires_val = requires_val


class Args:
    def __init__(self, usage: str):
        """
        Takes the arguments passed and parses them.
        It will create a dictionary from which we can fetch the values, if any.
        The keys are short.
        """
        self.flags: dict[str, ArgOpts] = dict()
        self.__osargs = sys.argv[1:]
        self.__usage_string = usage


    def add_argument(self, key: str, default: any, requires_val: bool, description: str):
        self.flags[key] = ArgOpts(default, requires_val, description)
        pass


    def __print_args(self):
        """
        Print the arguments in a nice way.
        """
        for k, v in self.flags.items():
            print(f' -{k} {"<value>\n" if v.requires_val else "\n"}{f"   {v.description}\n"}')


    def usage(self):
        """
        Prints the usage of the program along with the arguments.
        """
        print(self.__usage_string)
        self.__print_args()


    def parse_args(self):
        """
        This method will try to parse the passed arguments.
        In case of an error, the error will be raised, so surround this method with try/catch
        """
        # We go through all the 

        i = 0
        argslen = len(self.__osargs)

        while i < argslen:
            # We store the current argument
            sysarg = self.__osargs[i]

            # If it is a flag we know we parse it accordingly
            if sysarg[0] == '-' and sysarg[1] in self.flags.keys():
                # If the flag requires a value, we take the next argument
                if self.flags[sysarg[1]].requires_val:

                    # If the index jumps over the arguments length it means we did not receive a value
                    if i+1 >= argslen:
                        raise Exception(f'{sysarg} requires a value and it did not receive one')

                    # If the value is another flag it means the value for the current flag is missing
                    if self.__osargs[i+1][0] == '-' and self.__osargs[i+1][1] in self.flags.keys():
                        raise Exception(f'{sysarg} requires a value and it did not receive one')

                    # Otherwise we simply get it
                    self.flags[sysarg[1]].value = self.__osargs[i+1]
                    i += 1

                # If the flag does not require a value, we assume it is a bool, 
                # thus we turn it to the opposite of the default value
                else:
                    self.flags[sysarg[1]].value = not bool(self.flags[sysarg[1]].value)
                
                i += 1
            # Otherwise we raise an error
            else:
                raise Exception(f'unknown argument: {sysarg}')


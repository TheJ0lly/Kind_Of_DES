import args
import os

# EXIT CODES
OK = 0
HELP_CALLED = 1
ARGS_PARSING_ERROR = 2


if __name__ == "__main__":
    a = args.Args("SexyDES [OPTIONS]")
    a.add_argument('k', '', True, "the key used for encryption/decryption")
    a.add_argument('t', '', True, "the text to be encrypted/decrypted")
    a.add_argument('d', False, False, "starts the decryption process")
    a.add_argument('e', False, False, "starts the encryption process")
    a.add_argument('h', False, False, 'shows the help menu')


    try:
        a.parse_args()
    except Exception as e:
        print(f'error: {e}')
        os._exit(ARGS_PARSING_ERROR)

    if a.flags['h'].value == True:
        a.usage()
        os._exit(HELP_CALLED)
import args
import os
import bitset

# EXIT CODES
OK = 0
HELP_CALLED = 1
ARGS_PARSING_ERROR = 2
BOTH_CALLED = 3
NONE_CALLED = 4
WRONG_KEY_LEN = 5

if __name__ == "__main__":
    # We define the arguments
    a = args.Args("SexyDES [OPTIONS]")
    a.add_argument('k', '', True, "the key used for encryption/decryption")
    a.add_argument('t', '', True, "the text to be encrypted/decrypted")
    a.add_argument('d', False, False, "starts the decryption process")
    a.add_argument('e', False, False, "starts the encryption process")
    a.add_argument('h', False, False, 'shows the help menu')

    # We parse the arguments
    try:
        a.parse_args()
    except Exception as e:
        print(f'error: {e}')
        os._exit(ARGS_PARSING_ERROR)

    # If help has been toggled we print the usage and help
    if a.flags['h'].value == True:
        a.usage()
        os._exit(HELP_CALLED)

    if a.flags['d'].value and a.flags['e'].value:
        print('error: both the encryption and decryption flags have been toggled - must have only 1')
        os._exit(BOTH_CALLED)
    
    if not a.flags['d'].value and not a.flags['e'].value:
        print('error: neither encryption nor decryption flag has been toggled - must have 1')
        os._exit(NONE_CALLED)


    # We start to process the key and text
    key = bitset.Bitset(a.flags['k'].value)

    if key.len_bytes() != 8:
        print('error: key must be exactly 8 characters')
        os._exit(WRONG_KEY_LEN)
    
    text = bitset.Bitset(a.flags['t'].value)
    text.adjust_len(8)

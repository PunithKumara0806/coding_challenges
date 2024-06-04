#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <getopt.h>

#include "../lib/error_functions.h"

#define MAX_BUF 1024

const int B_COUNT = 1;
const int C_COUNT = 2;
const int W_COUNT = 4;
const int L_COUNT = 8;
const int V_FLAG = 16;

//process the file with option flags.
void process_file(int fd,int flag_opts){
    int i,count,inword;
    i = count = inword = 0;
    char buf[MAX_BUF];
    //count of line, byte, word and char
    int ln,bn,wn,cn;
    ln = bn = wn = cn = 0;

    while((count = read(fd,buf,MAX_BUF)) > 0){
        i = 0;
        for (;i < count;i++){
            if (buf[i] == '\n') 
                ln++;
            //count non-whitespace and whitespace borders
            if(buf[i] != ' ' && buf[i] != '\n' && buf[i] != '\t' && buf[i] != '\r'){
                inword = 1;
            }
            else{
                if(inword) 
                    wn++;
                inword = 0;
            }
            // use the continuation byte in utf-8 to see if byte is starting a code point or belongs to a code point
            if ((buf[i]&0xc0) != 0x80)
                cn++;
        }
        bn += count;
    }
    if (count == -1) {
        errExit("read failed");
    }
    if(inword) 
        wn++;

    if(flag_opts & L_COUNT){
        if(flag_opts & V_FLAG)
            printf("line : ");
        printf("%d ",ln);
    }
    if(flag_opts & W_COUNT){
        if(flag_opts & V_FLAG)
            printf("word : ");
        printf("%d ",wn);
    }
    if(flag_opts & C_COUNT){
        if(flag_opts & V_FLAG)
            printf("char : ");
        printf("%d ",cn);
    }
    if(flag_opts & B_COUNT){
        if(flag_opts & V_FLAG)
            printf("byte : ");
        printf("%d ",bn);
    }
}

void print_help_info(){
    char* info [] = {
        "custom wc [OPTIONS] [FILE]\n",
        "-h or --help option for displaying this message\n",
        "-b or --byte byte count\n",
        "-w or --word word count\n",
        "-c or --char char count\n",
        "-l or --line newline count\n",
        "-v or --display 'word', 'line' ... infront of their counts\n"
    };
    int i = 0,n = sizeof(info)/sizeof(char*);
    for( i = 0; i < n; i++){
        printf(info[i]);
    }
}

int main (int argc, char** argv){
    //initialize variables
    static struct option long_options []= {
        {"byte",no_argument,0,'b'},
        {"char",no_argument,0,'c'},
        {"word",no_argument,0,'w'},
        {"line",no_argument,0,'l'},
        {"help",no_argument,0,'h'},
        //used to display along with respective label
        {"verbose",no_argument,0,'v'}
    };
   //processing options
    int c,option_index;
    c = option_index = 0;
    int flag_opts = 0;
    while(1) {
        c = getopt_long(argc,argv,"cblwv",long_options,&option_index);
        if(c == -1) {
            break;
        }
        switch (c) {
            case 'c':
                flag_opts |= C_COUNT;
                // printf("known option : char\n");
            break;
            case 'b':
                flag_opts |= B_COUNT;
                // printf("known option : byte\n");
            break;
            case 'w':
                flag_opts |= W_COUNT;
                // printf("known option : word\n");
            break;
            case 'l':
                flag_opts |= L_COUNT;
                // printf("known option : line\n");
            break;
            case 'v':
                flag_opts |= V_FLAG;
            break;
            case 'h':
            default:
                print_help_info();
                return 1;
        }
    }
    //default flag, enable for all 
    if(flag_opts == 0){
        flag_opts |= C_COUNT;
        flag_opts |= B_COUNT;
        flag_opts |= W_COUNT;
        flag_opts |= L_COUNT;
    }
    int fd = 0;
    //process from stdin when no file is given
    if(optind >= argc){
        process_file(STDIN_FILENO,flag_opts);
        printf("\n");
    }

    //process non option (file names)
    int index;
    for(index = optind;index < argc;index++){
        // printf("non-options : %s\n",argv[index]);
        if( (fd = open(argv[index],0,O_RDONLY)) != -1){
            process_file(fd,flag_opts);
            printf(" %s\n",argv[index]);
            close(fd);
        }
        else
            errMsg("Could not open file '%s'\n",argv[index]);
    }
   //processing fil    // fputs(buf,stdout);
}


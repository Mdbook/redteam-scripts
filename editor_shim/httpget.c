#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

struct string {
  char *ptr;
  size_t len;
};

void init_string(struct string *s) {
  s->len = 0;
  s->ptr = malloc(s->len+1);
  if (s->ptr == NULL) {
    fprintf(stderr, "malloc() failed\n");
    exit(EXIT_FAILURE);
  }
  s->ptr[0] = '\0';
}

size_t writefunc(void *ptr, size_t size, size_t nmemb, struct string *s) {
  size_t new_len = s->len + size*nmemb;
  s->ptr = realloc(s->ptr, new_len+1);
  if (s->ptr == NULL) {
    fprintf(stderr, "realloc() failed\n");
    exit(EXIT_FAILURE);
  }
  memcpy(s->ptr+s->len, ptr, size*nmemb);
  s->ptr[new_len] = '\0';
  s->len = new_len;

  return size*nmemb;
}

const char * restrict getIP(int num) {
  CURL *curl;
  CURLcode res;
  char str[50];
  curl = curl_easy_init();
  if(curl) {
    struct string s;
    init_string(&s);
    if (num == 0){
      curl_easy_setopt(curl, CURLOPT_URL, "mdbook.me/ip.txt");
    } else if (num == 1){
      curl_easy_setopt(curl, CURLOPT_URL, "129.21.141.218/ip.txt");
    } else {
      strcpy(str, "10.100.0.101");
      free(s.ptr);
      char *ret = str;
      return ret;
    }
    
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writefunc);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &s);
    res = curl_easy_perform(curl);
    /* Check for errors */ 
    if(res != CURLE_OK){
      num++;
      free(s.ptr);
      curl_easy_cleanup(curl);
      return getIP(num);
    }
    strcpy(str, s.ptr);
    str[strcspn(str, "\n")] = 0;
    printf("%s\n", str);
    free(s.ptr);

    /* always cleanup */
    curl_easy_cleanup(curl);
  }
  char *ret = str;
  return ret;
}

// int main(){

// }
package article

const sampleArticleContentMentions = `{
   "type":"doc",
   "content":[
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"ARTICLE_ID",
                  "type":"article",
                  "title":"Article",
                  "time":"2019-08-19T19:15:00.534Z"
               }
            }
         ]
      },
	  {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"ARTICLE_NOT_FOUND_ID",
                  "type":"article",
                  "title":"Article not found",
                  "time":"2019-08-19T19:15:00.534Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"ARTICLE_NO_ACCESS_ID",
                  "type":"article",
                  "title":"Article no access",
                  "time":"2019-08-19T19:15:00.534Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"GROUP_ID",
                  "type":"group",
                  "title":"Group",
                  "time":"2019-08-19T18:41:19.673Z"
               }
            }
         ]
      },
	  {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"GROUP_NOT_FOUND_ID",
                  "type":"group",
                  "title":"Group not found",
                  "time":"2019-08-19T18:41:19.673Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"LIST_ID",
                  "type":"list",
                  "title":"List",
                  "time":"2019-08-19T18:41:19.673Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"LIST_NOT_FOUND_ID",
                  "type":"list",
                  "title":"List not found",
                  "time":"2019-08-19T18:41:19.673Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"LIST_NO_ACCESS_ID",
                  "type":"list",
                  "title":"List no access",
                  "time":"2019-08-19T18:41:19.673Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"TAG_NAME",
                  "type":"tag",
                  "title":"Tag",
                  "time":"2019-08-19T18:41:37.572Z"
               }
            }
         ]
      },
	  {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"TAG_NOT_FOUND_NAME",
                  "type":"tag",
                  "title":"Tag not found",
                  "time":"2019-08-19T18:41:37.572Z"
               }
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"USER_NAME",
                  "type":"user",
                  "title":"User",
                  "time":"2019-08-19T18:41:39.697Z"
               }
            }
         ]
      },
	  {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"mention",
               "attrs":{  
                  "id":"USER_NOT_FOUND_NAME",
                  "type":"user",
                  "title":"User not found",
                  "time":"2019-08-19T18:41:39.697Z"
               }
            }
         ]
      }
   ]
}`

// sample content to test backend rendering
const sampleArticleContent = `{  
   "type":"doc",
   "content":[  
      {  
         "type":"headline",
         "attrs":{  
            "level":2
         },
         "content":[  
            {  
               "type":"text",
               "text":"a"
            }
         ]
      },
      {  
         "type":"headline",
         "attrs":{  
            "level":3
         },
         "content":[  
            {  
               "type":"text",
               "text":"b"
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"text",
               "text":"c"
            },
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"bold"
                  }
               ],
               "text":"d "
            },
            {  
               "type":"mention",
               "attrs":{  
                  "id":"YRGaqMoa8w",
                  "type":"article",
                  "title":"Einfacher Artikel",
                  "time":"2019-08-17T22:47:06.436Z"
               }
            },
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"bold"
                  }
               ],
               "text":" "
            },
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"italic"
                  }
               ],
               "text":"e"
            },
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"strikethrough"
                  }
               ],
               "text":"f"
            },
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"underlined"
                  }
               ],
               "text":"g"
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"text",
               "marks":[  
                  {  
                     "type":"link",
                     "attrs":{  
                        "href":"https://emvi.com/"
                     }
                  }
               ],
               "text":"h"
            }
         ]
      },
      {  
         "type":"bullet_list",
         "content":[  
            {  
               "type":"list_item",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"i "
                        },
                        {  
                           "type":"mention",
                           "attrs":{  
                              "id":"max",
                              "type":"user",
                              "title":"Max Mustermann",
                              "time":"2019-08-17T22:46:47.807Z"
                           }
                        }
                     ]
                  }
               ]
            },
            {  
               "type":"list_item",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"j"
                        }
                     ]
                  },
                  {  
                     "type":"bullet_list",
                     "content":[  
                        {  
                           "type":"list_item",
                           "content":[  
                              {  
                                 "type":"paragraph",
                                 "content":[  
                                    {  
                                       "type":"text",
                                       "text":"k"
                                    }
                                 ]
                              }
                           ]
                        },
                        {  
                           "type":"list_item",
                           "content":[  
                              {  
                                 "type":"paragraph",
                                 "content":[  
                                    {  
                                       "type":"text",
                                       "text":"l"
                                    }
                                 ]
                              }
                           ]
                        },
                        {  
                           "type":"list_item",
                           "content":[  
                              {  
                                 "type":"paragraph",
                                 "content":[  
                                    {  
                                       "type":"text",
                                       "text":"m"
                                    }
                                 ]
                              }
                           ]
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {  
         "type":"horizontal_rule"
      },
      {  
         "type":"image",
         "attrs":{  
            "src":"http://localhost:4003/api/v1/content/DoB9mwd3ZV.png"
         }
      },
      {  
         "type":"ordered_list",
         "attrs":{  
            "order":1
         },
         "content":[  
            {  
               "type":"list_item",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"n"
                        }
                     ]
                  }
               ]
            },
            {  
               "type":"list_item",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"o"
                        }
                     ]
                  }
               ]
            },
            {  
               "type":"list_item",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"p"
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {  
         "type":"blockquote",
         "content":[  
            {  
               "type":"paragraph",
               "content":[  
                  {  
                     "type":"text",
                     "text":"q"
                  }
               ]
            },
            {  
               "type":"blockquote",
               "content":[  
                  {  
                     "type":"paragraph",
                     "content":[  
                        {  
                           "type":"text",
                           "text":"r"
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {  
         "type":"paragraph",
         "content":[  
            {  
               "type":"file",
               "attrs":{  
                  "file":"http://localhost:4003/api/v1/content/FfbLp3jwy0",
                  "name":"12QnJATs51WK_49mb",
                  "size":"49.28 MB"
               }
            }
         ]
      }
   ]
}`
const exportSampleDoc = `{
   "type":"doc",
   "content":[
      {
         "type":"paragraph",
         "content":[
            {
               "type":"file",
               "attrs":{
                  "file":"http://localhost:4003/api/v1/content/j06rqfiflKwSRgmtw5li.txt",
                  "name":"test.txt",
                  "size":"20.94 kB"
               }
            },
            {
               "type":"file",
               "attrs":{
                  "file":"http://localhost:4003/api/v1/content/j06rqfiflKwSRgmtw5li.txt",
                  "name":"test.txt",
                  "size":"20.94 kB"
               }
            }
         ]
      },
      {
         "type":"image",
         "attrs":{
            "src":"http://localhost:4003/api/v1/content/hoADCzmBrtfmzM3LlzDJ.jpg"
         },
         "content":[
            {
               "type":"paragraph"
            }
         ]
      },
      {
         "type":"pdf",
         "attrs":{
            "src":"http://localhost:4003/api/v1/content/asdf.pdf"
         }
      },
      {
         "type":"paragraph"
      },
      {
         "type":"paragraph",
         "content":[
            {
               "type":"file",
               "attrs":{
                  "file":"http://localhost:4003/api/v1/content/6WQPzMlQMNluLuFVZk1V.go",
                  "name":"XMKILA7ZXlK.go",
                  "size":"122 bytes"
               }
            }
         ]
      },
      {
         "type":"image",
         "attrs":{
            "src":"http://localhost:4003/api/v1/content/BcnWqHddeaQj773mSCcR.png"
         },
         "content":[
            {
               "type":"paragraph"
            }
         ]
      },
      {
         "type":"pdf",
         "attrs":{
            "src":"http://localhost:4003/api/v1/content/BQyOlhF3RqEw9vHJNEBN.pdf"
         }
      },
      {
         "type":"paragraph"
      }
   ]
}`

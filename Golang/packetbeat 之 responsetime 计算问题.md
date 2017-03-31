# packetbeat 之 responsetime 计算问题

标签（空格分隔）： packetbeat

---

## 问题

使用 packetbeat 进行抓包分析的主要目的之一就是**确定 request 和 response 之间 latency 时间**；而 responsetime 值对应的就是这个时间，因此该值的计算方式，以及准确性在实际使用中非常重要；

## 表现

在基于 packetbeat 的分析报告中，经常可以看到如下输出内容

```
responsetime(193 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-29T08:52:01.965Z","beat":{"hostname":"sunfeideMacBook-Pro.local","name":"sunfeideMacBook-Pro.local","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":4,"client_ip":"10.0.58.183","client_port":26776,"client_proc":"","client_server":"","ip":"10.0.10.58","method":"SMEMBERS","port":7602,"proc":"","query":"SMEMBERS restaurant:326200:contents","redis":{"return_value":"[]"},"resource":"restaurant:326200:contents","responsetime":193,"server":"","status":"OK","type":"redis"}


responsetime(193 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-29T08:52:01.965Z","beat":{"hostname":"sunfeideMacBook-Pro.local","name":"sunfeideMacBook-Pro.local","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.58.183","client_port":26776,"client_proc":"","client_server":"","ip":"10.0.10.58","method":"SMEMBERS","port":7602,"proc":"","query":"SMEMBERS food_restaurant:2186309:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:2186309:contents","responsetime":193,"server":"","status":"OK","type":"redis"}
```

从中可以看到几个数值相同的地方：

- <1> 和 <2> 的 "@timestamp" 数值相同；
- <1> 和 <2> 的 "responsetime" 数值相同；
- <1> 和 <2> 的 "client_port" 和 "port" 相同；

不同点在于：

- "query" 的内容不同；


看起来似乎是在同一个时间发出了不同 query 请求，但在获取相应的时间时，得到了相同的值；其实这个是由于 pipeline 的原因，如下图所示：

![redis pipeline 对 packetbeat 分析的影响](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis%20pipeline%20%E5%AF%B9%20packetbeat%20%E5%88%86%E6%9E%90%E7%9A%84%E5%BD%B1%E5%93%8D.png "redis pipeline 对 packetbeat 分析的影响")

还存在另外一种情况：即 response 大到需要在 TCP 协议层面分包发送；

日志如下：

```
{"@timestamp":"2017-03-29T08:52:06.501Z","beat":{"hostname":"sunfeideMacBook-Pro.local","name":"sunfeideMacBook-Pro.local","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":13507,"client_ip":"10.0.58.183","client_port":26776,"client_proc":"","client_server":"","ip":"10.0.10.58","method":"GET","port":7602,"proc":"","query":"GET search_area_restaurant_id:wwmthp","redis":{"return_value":"(lp1\nI500232\naI1777849\naI619754\naI1403142\naI523501\naI2350975\naI151675482\naI1505143\naI1069044\naI614157\naI1406283\naI571267\naI2069336\naI151695752\naI147163644\naI2239901\naI2256112\naI451580\naI1217250\naI381736\naI528265\naI634537\naI767398\naI475421\naI1505027\naI400907\naI152107900\naI491118\naI416144\naI150970671\naI484445\naI150058323\naI1048628\naI317081\naI552229\naI519132\naI1962390\naI2262288\naI150899176\naI1509387\naI1551358\naI994255\naI151004498\naI324628\naI871239\naI2020650\naI2059480\naI921549\naI343032\naI2251490\naI192074\naI305984\naI1333436\naI1559413\naI151691503\naI578980\naI2359311\naI1069859\naI551534\naI688005\naI1090091\naI147502706\naI2060663\naI1189870\naI1810063\naI150162189\naI1053551\naI1881428\naI1265721\naI2349741\naI429957\naI401978\naI2113721\naI956012\naI150969522\naI216694\naI752204\naI1913280\naI408229\naI489046\naI1973669\naI150973128\naI907450\naI1858036\naI2047470\naI2310933\naI1420537\naI396459\naI147508072\naI314288\naI1467056\naI150897864\naI1500140\naI295317\naI1260075\naI202927\naI150101723\naI1051636\naI901363\naI1126673\naI1148076\naI959501\naI421094\naI453764\naI1012326\naI604775\naI602103\naI2132709\naI150160637\naI150851767\naI85240\naI525522\naI1147925\naI150883092\naI1897103\naI1048713\naI456792\naI2235222\naI396450\naI2224988\naI150972019\naI1519190\naI1425052\naI150855322\naI151642336\naI1915130\naI2164310\naI2174705\naI729864\naI671254\naI1291063\naI2028744\naI1107294\naI52027\naI556973\naI150992786\naI200796\naI1035913\naI648504\naI150052844\naI1972231\naI1904258\naI1304882\naI792645\naI150136030\naI2308724\naI2308727\naI648695\naI951796\naI753048\naI1536919\naI1469356\naI497237\naI2020502\naI676226\naI2195107\naI952901\naI2267803\naI1024375\naI2203394\naI1321801\naI1906374\naI558568\naI394262\naI1988030\naI1056855\naI168082\naI675100\naI2247914\naI800164\naI944599\naI1522307\naI151694252\naI652848\naI1866111\naI1866441\naI852980\naI1783720\naI1148100\naI283750\naI1148102\naI1001539\naI1823360\naI572147\naI1119549\naI384890\naI434975\naI2173503\naI446636\naI202961\naI901370\naI635638\naI2174591\naI2199680\naI1385546\naI1001902\naI1900579\naI497519\naI1885574\naI1028157\naI1397383\naI150879109\naI223911\naI1935449\naI1147851\naI376707\naI151704231\naI943972\naI1974773\naI1866233\naI886223\naI659015\naI1860886\naI1085306\naI1050226\naI143482230\naI448164\naI429877\naI2367972\naI147587877\naI983749\naI1181122\naI1304339\naI150967237\naI626696\naI479785\naI1936405\naI509876\naI1524715\naI151003408\naI1872271\naI2188925\naI1337903\naI1768761\naI2140255\naI142942425\naI2055240\naI1060960\naI150135720\naI777810\naI570124\naI2055160\naI693420\naI71235\naI1785398\naI1048273\naI402562\naI2026673\naI660297\naI2083100\naI143688\naI1435957\naI2246665\naI2061987\naI1924181\naI893305\naI614226\naI150879529\naI151682687\naI1393682\naI630763\naI797508\naI1875434\naI1838449\naI694343\naI1049983\naI1863797\naI1192097\naI150895097\naI685287\naI381696\naI147523922\naI951867\naI1476534\naI1277845\naI2147422\naI150140415\naI1908701\naI1880458\naI2330461\naI1438015\naI150068727\naI1148535\naI1371845\naI2216946\naI150162011\naI1115466\naI1824764\naI1098816\naI144028081\naI1784923\naI2231055\naI2004869\naI150893394\naI2242511\naI1050045\naI150135634\naI150039604\naI1491841\naI1469786\naI143157238\naI1418345\naI151690609\naI1148089\naI662320\naI1461795\naI772502\naI552200\naI2238979\naI150156167\naI812167\naI1781115\naI706669\naI2221035\naI150093078\naI620689\naI732740\naI713059\naI1534419\naI150852898\naI525515\naI643672\naI707182\naI1185116\naI1276811\naI143325050\naI1423488\naI1842657\naI1053962\naI329073\naI777222\naI1469234\naI1547042\naI150898598\naI390550\naI672301\naI893121\naI1322086\naI461795\naI948080\naI497499\naI729717\naI662143\naI1157222\naI1306601\naI1157224\naI2256341\naI283716\naI1047991\naI1289969\naI150136007\naI723735\naI417641\naI280016\naI996864\naI144036295\naI1424001\naI1398126\naI1143989\naI1524495\naI565435\naI1054112\naI1767947\naI1965142\naI2160880\naI341681\naI1395098\naI145354205\naI679137\naI2355156\naI1350456\naI552106\naI478414\naI150884935\naI720549\naI630986\naI360173\naI1800584\naI1880122\naI2004524\naI150007275\naI620424\naI1350075\naI150024231\naI2132708\naI867622\naI150140660\naI1420923\naI1003837\naI1964620\naI1820812\naI495095\naI1181000\naI800117\naI602382\naI221984\naI925799\naI152107221\naI498411\naI1325055\naI550321\naI2100070\naI58897\naI679554\naI143845525\naI1220700\naI1847426\naI1879887\naI419552\naI274958\naI2003755\naI1933054\naI1055687\naI629928\naI638178\naI2374015\naI150160469\naI1339305\naI130659\naI220632\naI1052448\naI471485\naI150160459\naI1885033\naI1421051\naI2195803\naI150029456\naI150092202\naI637098\naI1047981\naI2162072\naI464591\naI143256668\naI542556\naI536273\naI926974\naI788329\naI1881438\naI494837\naI1253964\naI150989212\naI451522\naI1284806\naI873353\naI2054650\naI2242292\naI276741\naI1417938\naI1396276\naI1443248\naI150079414\naI1066834\naI683853\naI626558\naI457793\naI428432\naI1001827\naI1464245\naI298016\naI380131\naI2160210\naI536458\naI2174707\naI2051611\naI150966081\naI1067887\naI150067266\naI608775\naI1098366\naI1837891\naI1827507\naI1381659\naI1149100\naI1821182\naI150879305\naI1152394\naI656282\naI496122\naI514804\naI533571\naI1123881\naI558671\naI150024444\naI755494\naI151662814\naI426455\naI1985884\naI471419\naI579043\naI2124227\naI1490239\naI456483\naI620809\naI1934159\naI1543154\naI150888013\naI1808925\naI713226\naI614201\naI1333292\naI2057721\naI1253144\naI360189\naI348638\naI151636426\naI1053108\naI799261\naI1324794\naI1842691\naI1022858\naI672220\naI448048\naI1541355\naI151683261\naI617650\naI2240294\naI2229025\naI451145\naI2240115\naI1806621\naI2184379\naI676192\naI658865\naI500572\naI676211\naI666345\naI1818646\naI922398\naI147052169\naI479322\naI2244452\naI1821377\naI2343223\naI1254061\naI396443\naI454991\naI1861471\naI1120789\naI1188818\naI610802\naI2004597\naI1853586\naI1112393\naI1299379\naI1370283\naI1813328\naI1078572\naI367661\naI150133914\naI447617\naI1077576\naI2054832\naI1032223\naI2102191\naI2253138\naI791040\naI1045991\naI948613\naI615745\naI2253132\naI151683465\naI151707518\naI380665\naI199861\naI1966008\naI2004054\naI150136220\naI483411\naI197235\naI150078347\naI915311\naI421109\naI2279256\naI150141084\naI1196413\naI1215552\naI2369100\naI1886360\naI921051\naI1097079\naI150007781\naI605442\naI1500050\naI2248761\naI907720\naI658246\naI610859\naI1813130\naI1288854\naI729814\naI150970802\naI151670783\naI2016618\naI1858932\naI789864\naI499203\naI150985266\naI1982418\naI2106702\naI637196\naI959500\naI1276176\naI1500557\naI151000765\naI57969\naI1364294\naI1075231\naI150073222\naI2175437\naI1845454\naI207160\naI955251\naI1138334\naI55705\naI202172\naI150089181\naI150139956\naI359096\naI617227\naI989813\naI1867526\naI621501\naI2348033\naI150162018\naI905604\naI2292926\naI151662658\naI152107844\naI328440\naI1821535\naI368639\naI1106025\naI677161\naI366731\naI2226896\naI531886\naI1048248\naI1049833\naI2357301\naI150150866\naI150099687\naI641518\naI1174347\naI1957091\naI151679723\naI2081725\naI150068838\naI257966\naI289441\naI1882009\naI218485\naI150024394\naI478977\naI1491754\naI2087363\naI2347699\naI721296\naI1316540\naI1443575\naI2218745\naI323250\naI150109759\naI1400617\naI2373502\naI394691\naI57686\naI1414771\naI150004574\naI150102397\naI1782540\naI151712258\naI2038578\naI676424\naI1993215\naI612635\naI1371027\naI1434385\naI1302470\naI1542764\naI151683300\naI1050053\naI1529323\naI1148099\naI717827\naI1448545\naI785482\naI2173945\naI570713\naI1025750\naI809994\naI1205771\naI784084\naI712907\naI1985263\naI512910\naI1173373\naI841288\naI273179\naI384885\naI600419\naI1399431\naI2226946\naI2191208\naI578678\naI676515\naI2238494\naI1879484\naI1402987\naI2012441\naI370658\naI418394\naI150004665\naI2266752\naI435620\naI1860900\naI413669\naI2308756\naI328663\naI384868\naI151683232\naI1869389\naI774433\naI297542\naI2046073\naI858550\naI759652\naI554559\naI1169765\naI1470289\naI2080256\naI2090113\naI497154\naI996872\naI150136155\naI615166\naI1842073\naI151679514\naI1293568\naI328215\naI1175437\naI749192\naI531905\naI150882586\naI151663533\naI167537\naI932042\naI2063648\naI145393219\naI1466974\naI150984297\naI1388025\naI532442\naI362636\naI884993\naI1867993\naI147549414\naI150129972\naI1867995\naI545072\naI216906\naI2189184\naI782273\naI1305176\naI1206399\naI455922\naI2233420\naI1002643\naI755143\naI737823\naI1332995\naI150892727\naI2167871\naI421918\naI1414027\naI147651946\naI1056793\naI1513903\naI1379788\naI994975\naI1152776\naI2039877\naI1557870\naI151683126\naI1317057\naI439482\naI35325\naI1322148\naI355783\naI1838127\naI1859018\naI626642\naI943778\naI208830\naI1366083\naI1475461\naI618566\naI1820765\naI1338247\naI1303853\naI935574\naI491523\naI2178325\naI1111881\naI1272411\naI1325140\naI1392442\naI649651\naI141579384\naI1525734\naI1403058\naI1465204\naI546832\naI2254955\naI1173177\naI1076866\naI2209499\naI331891\naI2188309\naI1925129\naI1305245\naI150134117\naI1785724\naI1944078\naI1946664\naI399100\naI670855\naI1176588\naI1942753\naI1794204\naI2227381\naI712981\naI2176077\naI569305\naI150980547\naI2193913\naI741997\naI1856591\naI638611\naI1162751\naI485564\naI1516439\naI1255933\naI1857106\naI1174011\naI2265903\naI151636628\naI2195613\naI2219470\naI662611\naI1372031\naI1247859\naI151652973\naI783700\naI1395397\naI2238353\naI622844\naI602088\naI2238450\naI1260114\naI1815571\naI1564527\naI809230\naI2133829\naI1270021\naI359044\naI1136642\naI151634049\naI330320\naI1825270\naI285120\naI267630\naI714188\naI566835\naI1966124\naI601700\naI1772909\naI150972132\naI998163\naI805082\naI898863\naI1869570\naI394814\naI613504\naI2041701\naI2001764\naI390516\naI1316955\naI150024601\naI2043598\naI258071\naI1914875\naI978998\naI1452829\naI719080\naI1404611\naI317955\naI2089165\naI223957\naI339261\naI519314\naI1851371\naI1858589\naI666339\naI1818802\naI523669\naI151663446\naI1533500\naI2179881\naI2377997\naI2298903\naI1054228\naI1041235\naI327366\naI1976317\naI339278\naI647800\naI1906311\naI813834\naI225214\naI2231023\naI618254\naI1465965\naI1857893\naI900666\naI1805900\naI2237497\naI150851519\naI1446111\naI1857898\naI898013\naI1497468\naI2199972\naI1933647\naI150141413\naI1926405\naI200768\naI150966752\naI1282212\naI857924\naI150980614\naI2172785\naI706749\naI1053603\naI177314\naI225075\naI1863568\naI1771211\naI150135633\naI626607\naI330335\naI1832935\naI279430\naI1333337\naI150966796\naI2005155\naI150888694\naI2050029\naI784911\naI1122413\naI150988760\naI320563\naI1394481\naI606281\naI57393\naI1468519\naI622743\naI612483\naI1148094\naI1857768\naI784648\naI1868928\naI944849\naI1856543\naI2204655\naI428818\naI456441\naI748116\naI1907084\naI2212366\naI1321855\naI609517\naI1329921\naI1215803\naI656643\naI1547623\naI357128\naI2244308\naI150165363\naI1165552\naI1050969\naI2106914\naI939705\naI2049688\naI1086084\naI1023165\naI1936128\naI1904726\naI965410\naI151001009\naI362675\naI920031\naI1421282\naI54570\naI258185\naI698196\naI150108068\naI1946519\naI150034097\naI1322070\naI2198525\naI2198526\naI2009080\naI147507115\naI54051\naI1945546\naI991733\naI150893426\naI636849\naI1339899\naI1856490\naI2189372\naI2090110\naI2090111\naI660961\naI787181\naI1949161\naI1829284\naI2056224\naI496540\naI773141\naI1309476\naI1889967\naI485248\naI150162602\naI480865\naI572349\naI339116\naI150136695\naI151714099\naI1829831\naI150971024\naI269169\naI2235781\naI2205327\naI1245152\naI238003\naI198287\naI855455\naI623930\naI1875202\naI2330255\naI917337\naI732241\naI1052018\naI2349708\naI283259\naI443866\naI1248911\naI1876012\naI1261266\naI1863009\naI1343452\naI1874074\naI1244761\naI1938299\naI649981\naI613222\naI1068186\naI944756\naI612629\naI151683709\naI1937674\naI150034266\naI274892\naI611258\naI1842080\naI1426464\naI541439\naI638909\naI1433685\naI1317638\naI1148930\naI490118\naI1419329\naI1177524\naI2257506\naI2105111\naI152108886\naI1400714\naI1174452\naI1005571\naI899310\naI152104409\naI47448\naI196265\naI2164303\naI150966364\naI1542749\naI713875\naI200854\naI407144\naI2164305\naI1162156\naI2324055\naI2164309\naI2295121\naI151682794\naI2351043\naI203232\naI390423\naI190336\naI1110604\naI465501\naI2026146\naI1959486\naI1354504\naI1053950\naI658906\naI2011321\naI1477804\naI2197602\naI2236115\naI2298036\naI878545\naI2178983\naI2086242\naI1220802\naI2356896\naI152108081\naI151707402\naI230068\naI1181241\naI1380481\naI150032823\naI473258\naI405156\naI403191\naI524847\naI328788\naI150971061\naI1915322\naI547008\naI867583\naI1868637\naI1548366\naI150140102\naI996318\naI634682\naI210867\naI1514239\naI1418155\naI1857605\naI505297\naI1962977\naI1254704\naI2073935\naI2226523\naI289075\naI151667185\naI2282406\naI1106490\naI1902768\naI379118\naI674623\naI614070\naI433844\naI150033684\naI1535551\naI728631\naI1138833\naI1307109\naI787553\naI463688\naI2015953\naI150875420\naI150873232\naI150968591\naI1395974\naI712563\naI2188039\naI1813663\naI1856766\naI1913141\naI1469190\naI797620\naI366960\naI2109018\naI1943583\naI2102899\naI1383634\naI150970023\naI1015789\naI1184390\naI897598\naI59161\naI952004\naI1332569\naI1158370\naI2026064\naI1413104\naI150135346\naI150982721\naI732833\naI1370455\naI2037717\naI2043718\naI435226\naI691659\naI2354555\naI2354554\naI1966236\naI1869965\naI677799\naI429777\naI1139016\naI977382\naI332043\naI475394\naI1077561\naI1805940\naI1942808\naI1978120\naI1047990\naI1190998\naI150139931\naI460184\naI533970\naI2234331\naI329184\naI150973110\naI1381733\naI147488840\naI409588\naI150973116\naI2086570\naI2191636\naI1062603\naI380129\naI1946614\naI229546\naI2178332\naI355922\naI150089063\naI1180129\naI1404468\naI150985390\naI150149768\naI602099\naI1107362\naI2217805\naI1897139\naI150112867\naI1491212\naI535619\naI1148091\naI463418\naI1888342\naI2010299\naI1215656\naI1315160\naI2343038\naI150140135\naI1166089\naI1280165\naI564822\naI1772911\naI672616\naI1443806\naI1517057\naI621559\naI1516445\naI376757\naI622188\naI794910\naI341425\naI1464453\naI1152680\naI1561181\naI2189680\naI1016760\naI515629\naI1996001\naI1350972\naI2308759\naI2226916\naI727937\naI510318\naI1002441\naI390505\naI1768488\naI1398064\naI2072954\naI1189578\naI563770\naI2111156\naI1457398\naI150972260\naI150975571\naI2157094\naI1517461\naI1439153\naI150989348\naI1426962\naI724826\naI251140\naI1390352\naI1131843\naI676173\naI150989343\naI250395\naI569255\naI1419796\naI151634810\naI150966846\naI398792\naI2246706\naI1930974\naI1821394\naI485302\naI547068\naI785830\naI449683\naI1196943\naI1330514\naI1945109\naI408240\naI1053952\naI633621\naI2114207\naI985237\naI1051644\naI1262175\naI1191641\naI1857737\naI150062089\naI500613\naI859798\naI366888\naI1148095\naI147505499\naI2226562\naI908486\naI1560595\naI2218291\naI1219422\naI666360\naI699300\naI1395870\naI2366838\naI481986\naI2366835\naI1047996\naI150843500\naI521293\naI1470026\naI1179539\naI2164307\naI2231689\naI390223\naI2010243\naI1860783\naI1882430\naI2081663\naI2183166\na."},"resource":"search_area_restaurant_id:wwmthp","responsetime":31,"server":"","status":"OK","type":"redis"}
```

可以看到：

- 应用层协议：request 占了 52 字节（"bytes_in":52），response 占了 13507 字节（"bytes_out":13507）；
- "responsetime" 仅为 31ms ；

抓包截图如下：

![redis response 需要 TCP 分包发送时对 packetbeat 分析的影响](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis%20response%20%E9%9C%80%E8%A6%81%20TCP%20%E5%88%86%E5%8C%85%E5%8F%91%E9%80%81%E6%97%B6%E5%AF%B9%20packetbeat%20%E5%88%86%E6%9E%90%E7%9A%84%E5%BD%B1%E5%93%8D.png "redis response 需要 TCP 分包发送时对 packetbeat 分析的影响")

由此可知：

- packetbeat 计算 "responsetime" 时是基于首个 TCP segment 计算的（在该抓包中，response 的最后一个 segment 和 request 之间的时间差为 54ms）；
- "bytes_in" 和 "bytes_out" 的计算则是计算的全部 TCP segment 总和（应用协议数据长度）；


## 深入

- 确定产生上述日志的代码

```golang
func (redis *redisPlugin) newTransaction(requ, resp *redisMessage) common.MapStr {
    ...
	responseTime := int32(resp.ts.Sub(requ.ts).Nanoseconds() / 1e3)

	event := common.MapStr{
		"@timestamp":   common.Time(requ.ts),  // requ.ts 即 request 包的时间戳
		"type":         "redis",
		"status":       error,
		"responsetime": responseTime, // response 包的时间戳 - request 包的时间戳
		"redis":        returnValue,
		"method":       common.NetString(bytes.ToUpper(requ.method)),
		"resource":     requ.path,
		"query":        requ.message,
		"bytes_in":     uint64(requ.size),
		"bytes_out":    uint64(resp.size),
		"src":          src,
		"dst":          dst,
	}
    ...
}
```

- request 和 response 的关联方式

```golang
// 将 request 和 response 进行关联
func (redis *redisPlugin) correlate(conn *redisConnectionData) {
    ...
    // merge requests with responses into transactions
	for !conn.responses.empty() && !conn.requests.empty() {
		requ := conn.requests.pop()
		resp := conn.responses.pop()

		if redis.results != nil {
			// 将匹配的 request 和 response 封装成 transaction event
			event := redis.newTransaction(requ, resp)
			// 将 transaction event 发布出去（比如写出到 stdout 或 file 等）
			redis.results.PublishTransaction(event)
		}
	} 
    ...
}
```

> 这里其实存在一个问题：上述实现中由于采用的是 FIFO 原则，因此必须确保 request 和 response 是正确匹配的，但在 redis 协议的 master 和 slave 之间进行命令同步相关协议通信时，不满足该要求，会导致匹配错误产生；

- 将 redis 协议中的 request 和 response 分类保存和匹配

```golang
func (redis *redisPlugin) handleRedis(
	conn *redisConnectionData,
	m *redisMessage,
	tcptuple *common.TCPTuple,
	dir uint8,
) {
    ...

	if m.isRequest {
		conn.requests.append(m) // wait for response
	} else {
		conn.responses.append(m)
		redis.correlate(conn)
	}
}
```

- 设置 request 或 response 时间戳的地方

```golang
func (redis *redisPlugin) doParse(
	conn *redisConnectionData,
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
) *redisConnectionData {

	// 基于当前 tcp 流的数据流向进行选择处理
	st := conn.streams[dir]
	if st == nil {
		// 新建 stream 时，为当前 redis 消息指定时间戳，即新 request 或 response 
		st = newStream(pkt.Ts, tcptuple)   -- 1
		conn.streams[dir] = st
		if isDebug {
			debugf("new stream: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
		}
	}
	...
	for st.Buf.Len() > 0 {
		if st.parser.message == nil {
			// 为新的一条 redis 消息指定时间戳，即新 request 或 response
			st.parser.message = newMessage(pkt.Ts)  -- 2
		}
   ...
   redis.handleRedis(conn, msg, tcptuple, dir)
   ...
}
```

> 可以看到时间戳的指定是在收到 redis 协议 request 或 response 的第一个 pkt 时进行的；

- TCP 协议包到应用层协议包的处理

```golang
// 将 tcp 协议层面的数据包添加到
func (stream *TCPStream) addPacket(pkt *protos.Packet, tcphdr *layers.TCP) {
    ...
	if len(pkt.Payload) > 0 {
		// 调用各协议模块定义的 Parse 函数
		conn.data = mod.Parse(pkt, &conn.tcptuple, stream.dir, conn.data)
	}
    ...
}
```

- 将 pkt 分配到不同的 TCP stream

```golang
func (tcp *TCP) Process(id *flows.FlowID, tcphdr *layers.TCP, pkt *protos.Packet) {
    ...
	// 基于 pkt 确定 TCP stream
	stream, created := tcp.getStream(pkt)
	...
	stream.addPacket(pkt, tcphdr)
}
```

- 针对基于 gopacket 从底层收上来的 packet 进行协议的逐层解析

```golang
func (d *Decoder) OnPacket(data []byte, ci *gopacket.CaptureInfo) {
    ...
	// Ethernet 层给出的类型
	currentType := d.linkLayerType
	
	// NOTE: 这里就是为每个 packet 设置时间戳的位置
	//       可以看到时间戳实际上是来自 gopacket 中给的值
	packet := protos.Packet{Ts: ci.Timestamp}
    ...
    for len(data) > 0 {
        ...
		// 根据 packet 所属的 layerType 触发相应的回调函数(ICMP/UDP/TCP)
		processed, err = d.process(&packet, currentType)
        ...
    }
}
...
// 根据 packet 所属的 layerType 触发相应的回调函数(ICMP/UDP/TCP)
func (d *Decoder) process(
	packet *protos.Packet,
	layerType gopacket.LayerType,
) (bool, error) {
    ...
	case layers.LayerTypeTCP:
		debugf("TCP packet")
		d.onTCP(packet)
		return true, nil
	}
    ...   
}
...
func (d *Decoder) onTCP(packet *protos.Packet) {
    ...
	d.tcpProc.Process(id, &d.tcp, packet)
}
```

- sniffer 根据配置 DataSource 读取数据包

```golang
func (sniffer *SnifferSetup) Run() error {
    ...
	for sniffer.isAlive {
	    ...
		// 从指定数据源（live interface 或者 pcap 文件）中获取下一个存在的 packet 
		//  data:  The bytes of an individual packet.
	    //  ci:  Metadata about the capture
		data, ci, err := sniffer.DataSource.ReadPacketData()
		...
		sniffer.worker.OnPacket(data, &ci)
	}
    ...
}
```

> 由上下文可知，时间戳是从 ci 中得到的，因此 DataSource 使用哪种是关键；

- 数据源 DataSource 的选择

```golang
func (sniffer *SnifferSetup) setFromConfig(config *config.InterfacesConfig) error {
    ...
	switch sniffer.config.Type {
	case "pcap":
	    ...
	    sniffer.DataSource = gopacket.PacketDataSource(sniffer.pcapHandle)
	case "af_packet":
	    ...
	    sniffer.DataSource = gopacket.PacketDataSource(sniffer.afpacketHandle)
	case "pfring", "pf_ring":
	    ...
	    sniffer.DataSource = gopacket.PacketDataSource(sniffer.pfringHandle)
	default:
		return fmt.Errorf("Unknown sniffer type: %s", sniffer.config.Type)
	}
	...
}
```

> 对应测试中的实际情况而言，此处只考虑 pcap 对应的情况；

- ci 中时间戳的设置

在 `pcap.go` 中有

```golang
func (p *Handle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
    ...
	err = p.getNextBufPtrLocked(&ci)
    ...
}
...
func (p *Handle) getNextBufPtrLocked(ci *gopacket.CaptureInfo) error {
    ...
    // 设置 ci 中时间戳的位置
	ci.Timestamp = time.Unix(int64(p.pkthdr.ts.tv_sec),
		int64(p.pkthdr.ts.tv_usec)*1000) // convert micros to nanos
	ci.CaptureLength = int(p.pkthdr.caplen)
	ci.Length = int(p.pkthdr.len)
    ...
}
```

可以看到，时间戳是从 `p.pkthdr.ts.tv_sec` 和 `p.pkthdr.ts.tv_usec` 两个值得到的，而 pkthdr 的定义如下：

```golang
type Handle struct {
	// cptr is the handle for the actual pcap C object.
	cptr         *C.pcap_t
	blockForever bool
	device       string
	mu           sync.Mutex
	// Since pointers to these objects are passed into a C function, if
	// they're declared locally then the Go compiler thinks they may have
	// escaped into C-land, so it allocates them on the heap.  This causes a
	// huge memory hit, so to handle that we store them here instead.
	pkthdr  *C.struct_pcap_pkthdr
	buf_ptr *C.u_char
}
```

由此可知，要想知道 `pkthdr.ts` 的值是如何得到的，则需要研究 `libpcap` 代码了；


## 展开

在《[gettimeofday() should never be used to measure time](https://blog.habets.se/2010/09/gettimeofday-should-never-be-used-to-measure-time.html)》中，有如下说明：


> The `pcap_pkthdr` struct (the “received packet” struct) contains a `struct timeval ts` that ruins our ability to measure the time it takes for the reply you get for some query you sent. They tell me the kernel supplies the timestamp, so it’s not really libpcaps fault.
> 
> Calling `clock_gettime()` when `libpcap` gives you a packet has turned out to be useless, as **the time difference between packet reception and the delivery to your program is too long and unstable.** You’re stuck with this wall-clock time until you fix all the kernels in the world and break binary compatibility with old `libpcap` programs.


在《[how to understand the capture time!](https://lists.gt.net/ethereal/users/4665)》中，有如下精彩讨论：

- [Ask #1]

> Hi, Dear All, 
>
> The webpages about pcap says that the "`pcap_pkthdr`" structure contains the information about when the packet was sniffed, that is: 
>
> ```
> struct pcap_pkthdr{ 
>     struct timeval ts;
>     bpf_u_int32 caplen;
>     bpf_u_int32 len;
> } 
> ```
>
> I wonder **whether the "`ts`" is just the time when the pcap captured the packet?** Whether Ethereal use this data for the time when a packet was captured? 
> 
> Ethereal display the captured packets like: 
> ```
> Frame 1 
> Arrival time: Jun 13, 2002 12:00:00.1234546789 
> ¡¡ ... 
> ```
> 
> How Ethereal gets this arrival time? from the `pcap_pkthdr` mentioned upper? the datum "123456789" come directly from the "`tv_usec`" part in the `timeval` strcuture? 

- [Answer #1]

> What "`ts`" means **depends on the operating system on which you're capturing packets**.
> 
> **On most operating systems**, it's the time at which the driver for the network interface gave the packet to the OS's packet capture mechanism; **On some OSes where the operating system doesn't itself put a timestamp on the packet**, it's the time at which the `libpcap` library read the packet from the OS kernel. 
> 
> I.e., the time isn't necessarily the time when the packet arrived on the machine running `tcpdump`/`Ethereal`/whatever sniffer program you're using - it may be a later time (although it probably won't be much later). 
>
> The answers to all these questions are "YES":
> 
> - Whether Ethereal use this data for the time when a packet was captured? 
> - How Ethereal gets this arrival time? from the `pcap_pkthdr` mentioned upper? 
> - the datum "123456789" come directly from the "`tv_usec`" part in the `timeval` strcuture?
>
> Note that not all OSes necessarily provide high-precision timestamps; they might, for example, provide timestamps with 1 millisecond or 10 
millisecond resolution.

- [Ask #2]

> However, I still don't understand the capture time Ethereal display. for example, when I capture the icmp packet produced by "ping host B" on host A, it shows the same capture time of echo request and echo reply, as the following: 
> 
> ```
> 1 0.000000 A B ICMP Echo(ping) request 
> Arrival Time: Jun 14,2002 12:00:00.123456789 
> ... 
> 2 0.000000 B A ICMP Echo(ping) reply 
> Arrival Time: Jun 14,2002 12:00:00.123456789 
> ...
> ```
> 
> I wonder **why the set of icmp packets arrive at the same time?** since A ping B, and B returns a echo reply, it shouldn't produce at the same time! 
> 
> More ever, I captured the "A ping B" echo request packet on host B, and I want to compute the transmit time for the packet.(A and B have been synchronized by NTP time server)
> 
> But 
> **transmit = "the arrvial time on host B" subtract "the time of the echo request produced on host A"** 
>
> the transmit seems much different from the "round-trip time"/2 displayed by "ping", I mean, it seems they are not in the same quantity scale. So I feel confused. Would you like to give me some suggestion? 

- [Answer #2]

> You're assuming they *did* arrive at the same time; that may not be the case. 
> 
> Perhaps the OS on which you're running the program on which you're capturing packets doesn't timestamp the packets with a sufficiently high-resolution time stamp, and doesn't try to give packets unique time stamps, either. 
>
> If, for example, the reply was received .1 milliseconds after the request was sent, but the timer the OS uses to time-stamp the packets has only a 1-millisecond resolution, the two packets might be given the same time stamp even though the request was sent at a different time from when the reply arrived. 
> 
> There's nothing Ethereal can do about that; it just displays the times `libpcap` gave it (or gave whatever program wrote the capture file). 
> 
> There's probably not much `libpcap` can do about that, either; it just gets the times that the OS provides. 


在《[[tcpdump-workers] Monotonic clock timestamp on packets
](http://www.mail-archive.com/tcpdump-workers@lists.tcpdump.org/msg05260.html)》中，提到

- [Ask #1]

> Has anyone looked into timestamping the captured packets using `clock_gettime(CLOCK_MONOTONIC, ...)`?
> 
> I'm thinking adding a `struct timespec` to `struct pcap_pkthdr` and filling that in addition to the `struct timeval`.
> 
> For a request-reply situation a monotonic clock is much more reliable than `gettimeofday()`.

- [Answer #1]

> `pcap_pkthdr` is in a file.  You cannot add *ANYTHING* to it without breaking compatibility; you'd have to introduce a new magic number.
> 
> BTW, note that if you call `clock_gettime()`, **there is *NO* guarantee that the time it returns has anything to do with the time the packe arrived**; it tells you the time when it's called, not the time when the packet arrived.
> 
> The only platforms on which `libpcap` uses `gettimeofday()` are:
>
> - **DLPI platforms** - where the DLPI module doesn't supply the time stamp (e.g., HP-UX);
> - **DOS**;
> - **Septel devices**;
> - USB capturing on Linux if you're not using the memory-mapped interface.
> 
> On all other platforms - i.e., on most of the platforms where `libpcap` is used - the time stamps are supplied to userland by the kernel, so if you want to use a different timer, you'll have to modify the kernel.  (*BSD, Mac OS X, Linux, Solaris, etc.)

- [Ask #2]

> Exactly. That's why I asked if anyone has taken a look at it. Because calling it from the application at `pcap_dispatch` time would be useless. Just like calling it from `libpcap` an arbitrary time too late would be useless.
> 
> So if the underlying systems don't provide a monotonic clock for packet arrival time then that's that.


结论：

- 在有些系统上，这个时间戳对应的是网卡驱动将数据交付 OS 内核的时间；
- 在另外一些系统上，这个时间戳对应的是负责 sniffer 的东东（例如 libpcap）从 OS 内核读取数据包的时间；
- 应该认为这个时间戳并非准确的，尤其是在系统压力比较大的情况下；


----------


## 其他

- [[tcpdump-workers] Receive timestamp](http://www.mail-archive.com/tcpdump-workers@lists.tcpdump.org/msg02528.html)
- [[tcpdump-workers] libpcap based timestamp in linux](http://www.mail-archive.com/tcpdump-workers@lists.tcpdump.org/msg01926.html)
- [[tcpdump-workers] libpcap and select problem](http://www.mail-archive.com/tcpdump-workers@lists.tcpdump.org/msg00791.html)
- [CINBAD investigation of different packet filters ](https://openlab-mu-internal.web.cern.ch/openlab-mu-internal/03_Documents/3_Technical_Documents/Technical_Reports/2008/CINBAD_Investigation_of_Different_Packet_Filters.pdf)
- [Libpcap tutorial](http://wiki.ucalgary.ca/page/Libpcap_tutorial)
- [mail-archive](https://www.mail-archive.com/)
- [gt.net](https://lists.gt.net/)




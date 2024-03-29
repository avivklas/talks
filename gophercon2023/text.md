שלום, אני אביב והיום אני הולך לדבר איתכם על
Hot path optimizations for latency-critical applications in GO
פשוט הייתי חייב למצוא כותרת יותר ארוכה מהשם שלי אחרת זה היה יוצא קצת מוזר

הדחף שלי לעמוד פה ולדבר על הנושא הזה בא בעיקר מתשוקה אבל גם מזה שאני מרגיש שהשפה הזו שאני כל כך אוהב קצת מפספסת אוהדים לאחרונה בגלל שהיא מנסה להתאים לכל משימה. זה נכון שהמטרה העיקרית שלה היא רמת פרודקטיביות גבוהה, ומכאן שאפשר לייצר בעצרת הכלים שבאים איתה תוכנה נהילה בקצב גבוה אך במחיר של ביצועים, בעיקר בגלל מודל הזכרון שנשען על GC, אבל זה לא אומר שלא ניתן לכתוב בעזרתה קוד אופטימלי מבחינת שימוש בזכרון שלא יבייש תוכנית הכתובה ב C מבחינת ביצועים.

לפני שנמשיך, אציג את עצמי בקצרה. שמי כאמור אביב, אני מפתח עם המון אהבה לתחום מגיל מאד צעיר וכיום מתפקד כארכיטקט מוביל בחברת Cyolo שם אנחנו מפתחים פתרון גישה מאובטחת למשאבים ארגוניים של הדור הבא (מן תחליף לVPN), בעזרת טכנולוגיה מקורית ומאד מעניינת מבחינה הנדסית.

מישהו פה לא בטוח שהוא איתי לגבי המונח HOT PATH?
אז אחדד. HOT PATH הוא איזור בקוד שרץ הרבה מאד פעמים ושיש לו השפעה על תפקוד המערכת מבחינה מוצרית. כלומר, איכות המוצר תפגע או תשתפר אם הוא יעבוד לאט או מהר יותר.

כדבר הראשון, הייתי רוצה לפתוח בכמה מילים על benchmarking
יש לי סיפורים מעניינים על הפתעות שגיליתי בעזרת benchmarking כולל הפסדים בהתערבויות שיספיקו להרצאה בפני עצמה, אז רק אומר שמאד חשוב לא לסמוך תמיד על האינטואיציה שלכם ופשוט מאד לבדוק. אם קשה לכם עם האגו, אתם תמיד יכולים לזרוק “לך תדע איזה אופטימיזציות סטטיות הוסיפו שם בקומפיילר לאחרונה” כהסבר ללמה אתם צריכים לערוך benchmark.
ספציפית בGO הכלי הזה כל כך פשוט ונגיש שזה באמת פשע לא להשתמש בו. לעבוד וליצור קוד חשוב שמבוסס על הנחות שוא יכול להגמר במקרה הטוב בתיקון אחרי עליה לאויר ובמקרה הגרוע לידע אבוד. ודבר אחרון בנושא, תמיד לדווח אלוקציות, כי אחרת איך נדע מה לשפר?!

עכשיו בחזרה לנקודה עליה דיברתי בהתלחה ומשהו שהוא די מובן מאליו לכל מי שכותב בGO - המבנה של חבילת IO שהוא בסיס לכל כך הרבה כלים בהם אנחנו משתמשים יום יום.
יש פה מישהו שלא מבין על מה אני מדבר?
תראו את החתימה של פונקציית Read.
הסיבה שהיא מקבלת סלייס של בתים היא כדי שהשליטה על היעד של המידע תהיה של הקורא לפונקציה.
כל קריאה ממשאב חיצוני משתמשת באותה חתימה וכך ניתן מאד בקלות לממש אלגוריתם אופטימלי מבחינת שימוש בזכרון ומכאן גם בזמן ריצה וחסכון בהשענות על ה GC.
דוגמא טובה לכך היא כל הכלים של חבילת bufio ו json decoder שמשמשים קריאה יעילה של משאבים.

עכשיו בטח חצי מהקהל שואל את עצמו על מה מדבר האיש המוזר הזה ולמה בכלל צריך להתעסק בזה. אז אנסה להסביר בעזרת ניסוי קטן. יצרתי קובץ גדול המכיל טקסט אנושי. הפסקאות בו המופרדות בירידת שורה כמובן. גודל הקובץ הוא 6M. כתבתי שני מימושים לספירת מופעים של מילה מסויימת בטקסט. אחד ע”י שימוש ב bufio.Reader שממש buffered reader. והשני הוא מימוש דומה, משלי. ההבדל הוא בעצם היכולת שלי לקרוא כל שורה לאותו באפר כי השתמשתי בפונקציה הנמוכה יותר שממשת io.Read. התוצאות מראות ששיפרתי ב25%. איך? פשוט מאד. לא היה צורך בהעתקת סלייסים כדי להחזיר תשובה מפונקציה.

זוהי בעצם רק דוגמא שבאה להוכיח את הנקודה שאפשר מצד אחד להינות מכלי השפה הזו ולכתוב מאד מהר, אבל כשצריך,
אפשר גם לעבוד טיפה יותר קשה ולכתוב קוד אופטימלי מבחינת שימוש בזכרון שבכלל לא מייצר הפניות ב GC.
בכל סיטואציה שתרצו תוכלו לנהל זכרון באופן שכזה ולממש לוגיקה ב hot path ללא אלוקציות.

הנושא הבא הוא סביב מקביליות. זה אמנם תענוג לכתוב קוד מקבילי בGO אך קורה הרבה שאנחנו מפספסים את מה שקורה באמת בזמן ריצה - צווארי בקבוק. אז הכנתי כמה דוגמאות לאיך אפשר להמנע מהם או איך למזער אותם

1. לפעמים אנחנו רוצים לעדכן משתנה פרימיטיבי באופן מקבילי ושוכחים שאפשר מאד בקלות לעשות זאת ללא מנעול
2. כשאין ברירה וצריך להשתמש במנעול - רצוי מאד שהוא יהיה הדדי. כלומר לאפשר מקביליות אמיתית בקריאה. אפשר כך לפתור כל מקרה, לפעמים ע”י שימוש בטכניקה המפורסמת - צ’ק-לוק-צ’ק - העיקר שהמקרה הנפוץ יותר ישתמש בנעילת קריאה בלבד.
3. בונוס: כמה פעמים יצא לכם (בפיתוח אמיתי, לא בראיונות עבודה) להפעיל פול של וורקרים שלוקחים עבודה מתור אחד במקביל? הרבה פעמים יוצא שהתוצאות של העבודה צריכות להסתנכרן בגלל גישה לזכרון משותף וכאן יש קצת אשליה. אין פה באמת מקבול כי הכל בסוף מסתנכרן. פתרון שאפשר לממש ללא הרבה מאמץ הוא כתיבה בבאצ’ים. כך:

הנושא האחרון לא קשור באופן ישיר לGO אך אני מרגיש שרבות הפעמים בהן אנו משלים את עצמנו לגבי מקבול.
אני רוצה לדון בהבדל בין
Concurrency
לבין
Parallelism

האחד הוא יכולת והשני הוא מעשה.
ביכולת מקבול הכוונה היא לפעולות שמאפשרות לכמה רוטינות להתייחס לאובייקטים משותפים.

אני מרגיש שהרבה פעמים אנחנו קצת מפספסים את המטרה
כשאני מפעיל כמה רוטינות שפועלות במקביל על אותה משימה, אני לא בהכרח ממקבל את הלוגיקה שלי וזאת מכמה סיבות:
- אם בסופה של עבודה נדרש כל פועל לעדכן את אותה כתובת בזכרון אזי שמתבצע פה סנכרון וכולם נאלצים לחכות באותו התור. זה קצת כמו בכביש החוף במחלף חבצלת איפה ששלושה נתיבים הופכים לשניים או אם בסופה של עבודה מתעדכן משאב חיצוני כמו מסד נתונים, אותו סנכרון תלוי במנגנון החלוקה של אותו מסד. אם מדובר בדטבייס רלציוני אין שום חלוקה כשמדובר על טרנזקציה שמעדכנת טבלה ושם נוצר מצב דומה למחלף חבצלת רק יותר גרוע כי הפעם מדובר בשרות שמשרת עבודות אחרות. לעיקרון הזה בתחום במקבול קוראים תלות מידע.
- דבר נוסף הוא שאם אני מתעלם מכמות המעבדים הזמינה לי, אני בעצם מייצר תחושה מדומה של מקבול. הרי אין כזה דבר באמת מקבול על ליבה חישובית אחת. בעצם קוראים לזה TASK SWITCHING. כל רוטינה תקבל את הכוח/כמות הפעימות היחסי לה מתוך סך הרוטינות

בואו נעבור לדוגמאות מקבול (לעומת יכולת מקבול) שאנחנו מכירים כדי לקבל קצת יותר תחושה של מה שאני מתכוון אליו:
- ספארק - מנוע להרצת פעולות חישוביות על זרמים של מידע, במקביל. ספארק מציע שפה שבעזרתה הוא יכול להבין אילו חלקים באפליקציה שלנו הוא יכול להריץ במקביל על יחידות חישוביות שונות ללא תלויות, ממש דומה לטכניקה המפורסמת mapReduce
- קאפקא - מציע טכניקה של צרכנים מודעים אותם הוא מנהל ובכך יכול לייצר חלוקה של התור לחלקים צפויים - יעני שארדינג - ולאפשר לנו לממש מקבול אמיתי בין מחשבים שונים - אם בסופו של קו הצינורות אנחנו לא פוגשים את מחלף חבצלת
- חלוקה במסדי נתונים - בהרבה מאד מסדי נתונים חלוקה לא רק באה ליעל רק את הכתיבות לדיסק אלא גם את הנעילות שנדרשות כדי לאפשר לצרכן לכתוב אליו מתהליכים מקביליים. גם כאן המנגנון מאחורי הקלעים יהיה שארדינג. כל שארד יהיה אחראי על טווח מסויים וינהל טור של עדכונים.

אז איך אנחנו, האנשים הקטנים, יכולים לממש מקבול:
- הכל מתחיל ונגמר בצינור המידע. עלינו ראשית לזהות האם המידע שעלינו לעבד יכול לעבור כל הדרך באופן מחולק בין יחידות העיבוד השונות שאחראיות על טווח שונה
- אז נוכל להפעיל את הפועלים על אותם טווחי חלוקה בהתאם. מכירים את זה שאתם באים לאסוף כרטיסים למופע והתורים מחולקים לפי טווחי אותיות בהן מתחיל שם המשפחה שלכם? אז בדיוק כך
- כפי שהזכרתי, זה יהיה חסר משמעות להפעיל שני עובדים על ליבה אחת ולכן הדבר הנכון יהיה להפעיל עובדים לפי כמות הליבות הזמינה לנו

לסיכום:
- דיברנו על נפלאות החבילה io וראינו שלפעמים שווה להתאמץ ולהשתמש בעקרונות שלה ישירות כדי להשיג שיפור ביצועים
- ראינו איך אנחנו יכולים לפתור צווארי בקבוק כשאין ביכולתנו למקבל
- למדנו איך לממש מקבול אמיתי מתי שאפשר

תודה רבה!

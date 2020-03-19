functionality:
    paivalehti init // later
    paivalehti render
    paivalehti add
    paivalehti social sync


paivalehti add <url> --description

    image: local, remote
    album: local, remote
    post: yaml?
    url: link

storage:

    <date>
        <type>.yaml
        <type>.jpg



    root/<date>/cover.jpg
    root/<data>/thumbnails/<resolution>-<hash>-cover.jpg


    root/<year>-<month>-<day>/tweet.yaml

        posts:
            - kind: tweet
            timestamp: x
            content: y
            avatar: z


    root/<year>-<month>-<day>/cover.yaml

        posts:
            - kind: photo
            ... 


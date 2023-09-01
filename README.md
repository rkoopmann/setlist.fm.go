# setlist.fm.go

## purpose

To gather a user's concert attendance listings from setlist.fm via api calls. Currently, this will write an overall list.json file with top-level event information andlong with individual event files following the naming convention YYYY-MM-DD-[artist name].json. This assumes the user has not attended multiple events by the same artist on the same day.

## motivation

Half learning, half self-serving. I wanted to learn golang and I also wanted this to exist in the world. so here we are.

## configuration

Refer to `sample-configuration.json` for configuration file layout.

- `user` -- username of user you'd like to stalk. Hopefully it's yourself.
- `apiKey` -- head over to https://api.setlist.fm/docs/1.0/index.html and apply for an api key; plunk that value here.
- `outputPath` -- file path for output json files. If you're using these for Jekyll, you might want to drop these into the `_data` directory -- see [the docs](https://jekyllrb.com/docs/datafiles/).

## building & running

```bash
# buildit
go build setlist_fm.go

# run it
./setlist_fm
```

## nuts & bolts

Built on, and run with `go version go1.13.5 darwin/amd64`. There are no test methods yet; still learning go.

## roadmap

- [ ] generate Apple music playlist of setlist for easy importing into Apple Music app.

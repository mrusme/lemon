#!/usr/bin/env python3
import falcon
import json
import ledhat

class MiddlewareJson(object):
    # Borrowed this from https://eshlox.net/2017/08/02/falcon-framework-json-middleware-loads-dumps/
    def process_request(self, req, resp):
        if req.content_length in (None, 0):
            return

        body = req.stream.read()

        if not body:
            raise falcon.HTTPBadRequest(
                'Empty request body. A valid JSON document is required.'
            )

        try:
            req.context['request'] = json.loads(body.decode('utf-8'))
        except (ValueError, UnicodeDecodeError):
            raise falcon.HTTPError(
                falcon.HTTP_753,
                'Malformed JSON. Could not decode the request body.'
                'The JSON was incorrect or not encoded as UTF-8.'
            )

    def process_response(self, req, resp, resource, req_succeeded):
        if 'response' not in resp.context:
            return

        resp.body = json.dumps(
            resp.context['response'],
            default=json_serializer
        )

class GitHubResource(object):
    def on_post(self, req, resp):
        if req.context['request']:
            body = req.context['request']
            print(body)

            repository_name = '(null)'
            repository_full_name = '(null)'
            if 'repository' in body:
                repository_name = body['repository']['name']
                repository_full_name = body['repository']['full_name']

            if 'head_commit' in body and 'id' in body['head_commit']:
                committer_username = '@' + body['head_commit']['committer']['username']
                commit_message = body['head_commit']['message']

                text = 'New commit by ' + committer_username + ' in ' + repository_full_name + ': ' + commit_message
                ledhat.icon('octocat')
                ledhat.text(text)
            elif 'forkee' in body and body['forkee']['fork'] == True:
                full_name = body['forkee']['full_name']
                text = repository_full_name + ' was forked to ' + full_name
                ledhat.icon('fork')
                ledhat.text(text)
            elif 'context' in body and body['context'] == 'ci/dockercloud':
                description = body['description']
                text = repository_full_name + ': ' + description
                ledhat.icon('docker')
                ledhat.text(text)

            resp.status = falcon.HTTP_204
        else:
            resp.status = falcon.HTTP_500

class IftttResource(object):
    def on_post(self, req, resp):
        if req.context['request']:
            body = req.context['request']
            print(body)
            # TODO: Find a way to integrate with IFTTT Webhooks (https://ifttt.com/maker_webhooks)
            # but still allow people to customize what would appear on lemon.
            resp.status = falcon.HTTP_204

app = falcon.API(middleware=[
    MiddlewareJson(),
])

github = GitHubResource()
ifttt = IftttResource()
app.add_route('/github', github)
app.add_route('/ifttt', ifttt)

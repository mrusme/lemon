#!/usr/bin/env python3
# coding=utf8

import falcon

class ResourceGitHub(object):
    def __init__(self, ledhat=None):
        self._ledhat = ledhat

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
                self._ledhat.text(text)
            elif 'forkee' in body and body['forkee']['fork'] == True:
                full_name = body['forkee']['full_name']
                text = repository_full_name + ' was forked to ' + full_name
                self._ledhat.icon('fork')
                self._ledhat.text(text)
            elif 'context' in body and body['context'] == 'ci/dockercloud':
                description = body['description']
                text = repository_full_name + ': ' + description
                self._ledhat.icon('docker')
                self._ledhat.text(text)

            resp.status = falcon.HTTP_204
        else:
            resp.status = falcon.HTTP_500

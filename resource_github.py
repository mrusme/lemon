#!/usr/bin/env python3
# coding=utf8

import falcon

class ResourceGitHub(object):
    def __init__(self, ledhat=None, influx=None):
        self._ledhat = ledhat
        self._influx = influx

    def on_post(self, req, resp):
        if req.context['request']:
            body = req.context['request']
            print(body)

            repository_name = '(null)'
            repository_full_name = '(null)'
            icon = None
            text = None
            category = 'undefined'

            if 'repository' in body:
                repository_name = body['repository']['name']
                repository_full_name = body['repository']['full_name']

            if 'head_commit' in body and 'id' in body['head_commit']:
                committer_username = '@' + body['head_commit']['committer']['username']
                commit_message = body['head_commit']['message']

                icon = 'octocat'
                text = 'New commit by ' + committer_username + ' in ' + repository_full_name + ': ' + commit_message
                category = 'github_commit'
            elif 'forkee' in body and body['forkee']['fork'] == True:
                full_name = body['forkee']['full_name']
                icon = 'fork'
                text = repository_full_name + ' was forked to ' + full_name
                category = 'github_fork'
            elif 'context' in body and body['context'] == 'ci/dockercloud':
                description = body['description']
                icon = 'docker'
                text = repository_full_name + ': ' + description
                category = 'github_ci_dockercloud'

            self._ledhat.icon('octocat')
            self._ledhat.text(text)

            if self._influx != None:
                self._influx.write(resource=req.path.replace('/','_'), icon=icon, category=category)

            resp.status = falcon.HTTP_204
        else:
            resp.status = falcon.HTTP_500

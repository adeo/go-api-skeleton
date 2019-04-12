#!/bin/bash
command -v gsed > /dev/null 2>&1 /dev/null && SED_CMD=gsed || SED_CMD=sed

OLD_PROJECT_NAMESPACE="github.com/adeo"
OLD_PROJECT_NAME="turbine-go-api-skeleton"
OLD_PROJECT_FULL_NAME="${OLD_PROJECT_NAMESPACE}/${OLD_PROJECT_NAME}"
OLD_PROJECT_GROUP="go-group-skeleton"

echo "What is the new project group? (Eg. turbine, qualite)"
read NEW_PROJECT_GROUP
echo "What is the new namespace? (Eg. if you are creating project under 'github.com/foo/bar', enter 'github.com/foo')"
read NEW_PROJECT_NAMESPACE
echo "What is the new project name? (Eg. if you are creating project under 'github.com/foo/bar', enter 'bar'))"
read NEW_PROJECT_NAME

NEW_PROJECT_FULL_NAME="${NEW_PROJECT_NAMESPACE}/${NEW_PROJECT_NAME}"
OLD_PROJECT_GIT="$(echo "$OLD_PROJECT_NAMESPACE" | $SED_CMD 's|/|:|')"
NEW_PROJECT_GIT="$(echo "$NEW_PROJECT_NAMESPACE" | $SED_CMD 's|/|:|')"

find . -iname '*.go' -exec $SED_CMD -i "s|${OLD_PROJECT_FULL_NAME}|${NEW_PROJECT_FULL_NAME}|g" {} \;
$SED_CMD -i "s|${OLD_PROJECT_FULL_NAME}|${NEW_PROJECT_FULL_NAME}|g" Makefile Dockerfile go.mod
$SED_CMD -i "s|${OLD_PROJECT_NAME}|${NEW_PROJECT_NAME}|g" Makefile Dockerfile info.yaml cmd/root.go tom.yml .gitlab-ci.yml
$SED_CMD -i "s|${OLD_PROJECT_GROUP}|${NEW_PROJECT_GROUP}|g" Makefile Dockerfile tom.yml .gitlab-ci.yml
$SED_CMD -i "s|${OLD_PROJECT_GIT}|${NEW_PROJECT_GIT}|g" .gitlab-ci.yml

main()
{
    echo
    echo
    echo "Now we will create entities and corresponding CRUD. Hit Ctrl+C to stop."
    echo
    while true
    do
        echo "Creating a new entity:"
        echo "What is the entity name you want to create? (name to be u$SED_CMD in URL path, write it lower case, singular)"
        read ENTITY_NAME
        ENTITY_NAME_UP=$(echo $ENTITY_NAME | tr  '[:upper:]' '[:lower:]')

        echo "What is the postgresql schema to use for this entity? (if you plan to use MongoDB you can ignore this question)"
        read ENTITY_SCHEMA

        cp handlers/template_handler.go handlers/${ENTITY_NAME}_handler.go
        $SED_CMD -i -r "s/template/${ENTITY_NAME}/g" handlers/${ENTITY_NAME}_handler.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" handlers/${ENTITY_NAME}_handler.go

        cp storage/dao/postgresql/database_postgresql_template.go storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/template/${ENTITY_NAME}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/schema/${ENTITY_SCHEMA}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go

        cp storage/dao/mongodb/database_mongodb_template.go storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/template/${ENTITY_NAME}/g" storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go

        cp storage/dao/mock/database_mock_template.go storage/dao/mock/database_mock_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/template/${ENTITY_NAME}/g" storage/dao/mock/database_mock_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/mock/database_mock_${ENTITY_NAME}.go

        cp storage/dao/fake/database_fake_template.go storage/dao/fake/database_fake_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/template/${ENTITY_NAME}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go

        cp storage/model/template.go storage/model/${ENTITY_NAME}.go
        $SED_CMD -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/model/${ENTITY_NAME}.go

        $SED_CMD -i -r "/\/\/ start: template routes/{:next;N;/\/\/ end: template routes/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" handlers/handler.go

        $SED_CMD -i -r "/\/\/ start: template dao funcs/{:next;N;/\/\/ end: template dao funcs/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" storage/dao/database.go

        $SED_CMD -i -r "/\/\/ template export/{p;s/Template/${ENTITY_NAME}/g}" storage/dao/fake/database_fake.go

        echo "Done"
        echo "If you want to stop here, hit Ctrl+C"
        echo
        echo
    done
}

main

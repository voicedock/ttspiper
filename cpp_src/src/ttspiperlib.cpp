#include "ttssimplelib/ttssimplelib.h"
#include "piper.hpp"
#include <vector>
#include <string>
#include <optional>

#ifdef __cplusplus
extern "C" {
#endif

typedef void (*textToAudioCb)(void);

void terminate()
{
    piper::terminate();
}

void initialize()
{
    auto exePath = filesystem::canonical("/proc/self/exe");
    piper::initialize(exePath.parent_path());
}

void textToAudio(void* voice, char *text, int cbId)
{
    void textToAudioCb(int cbId, int16_t *audioBuf, int audioBufLen);
    piper::SynthesisResult cppResult;

    std::vector<int16_t> buf;

    auto audioCallback = [cbId, &buf]() {
        int16_t audioBuf[buf.size()];

        std::copy(buf.begin(),
                  buf.end(),
                  audioBuf);
        textToAudioCb(cbId, &audioBuf[0], buf.size());
    };

    piper::Voice* vic = reinterpret_cast<piper::Voice*>(voice);


    piper::textToAudio(
        *vic,
        std::string (text),
        buf,
        cppResult,
        audioCallback
    );
//
//    result->inferSeconds = cppResult.inferSeconds;
//    result->audioSeconds = cppResult.audioSeconds;
//    result->realTimeFactor = cppResult.realTimeFactor;
}

void* loadVoice(char *modelPath, char *modelConfigPath, int64_t *speakerId) {
    piper::Voice *voice;
    voice = new piper::Voice;
    std::optional<piper::SpeakerId> optSpeakerId = std::nullopt;
    if (speakerId) {
        optSpeakerId = static_cast<int64_t>(*speakerId);
    }

    piper::loadVoice(modelPath, modelConfigPath, *voice, optSpeakerId);

    return voice;
}

#ifdef __cplusplus
}
#endif

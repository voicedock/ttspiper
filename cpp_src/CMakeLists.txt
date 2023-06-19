cmake_minimum_required(VERSION 3.14)
project(ttssimplelib
    VERSION 1.0.0
    DESCRIPTION "TTS simple lib"
    LANGUAGES CXX)
    
#----------------------------------------------------------------------------------------------------------------------
# general settings and options
#----------------------------------------------------------------------------------------------------------------------

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(GNUInstallDirs)

add_library(ttssimplelib SHARED
    src/ttssimplelib.cpp
)

set_target_properties(ttssimplelib PROPERTIES VERSION ${PROJECT_VERSION})
set_target_properties(ttssimplelib PROPERTIES SOVERSION 1)

#----------------------------------------------------------------------------------------------------------------------
# ttssimplelib dependencies
#----------------------------------------------------------------------------------------------------------------------
find_package(PkgConfig)
pkg_check_modules(ESPEAK_NG REQUIRED espeak-ng<2)
set(ONNXRUNTIME_ROOTDIR ${CMAKE_CURRENT_LIST_DIR}/lib/${CMAKE_HOST_SYSTEM_NAME}-${CMAKE_HOST_SYSTEM_PROCESSOR})
target_link_libraries(ttssimplelib
                      onnxruntime
                      pthread
                      -static-libgcc -static-libstdc++
                      ${ESPEAK_NG_LIBRARIES})

target_link_directories(ttssimplelib PRIVATE
                        ${ESPEAK_NG_LIBRARY_DIRS}
                        ${ONNXRUNTIME_ROOTDIR}/lib)

target_include_directories(ttssimplelib PRIVATE
                           ${ONNXRUNTIME_ROOTDIR}/include
                           ${ESPEAK_NG_INCLUDE_DIRS})

target_compile_options(ttssimplelib PUBLIC
                       ${ESPEAK_NG_CFLAGS_OTHER})

#----------------------------------------------------------------------------------------------------------------------
# ttssimplelib sources
#----------------------------------------------------------------------------------------------------------------------
include_directories(include)
set_target_properties(ttssimplelib PROPERTIES PUBLIC_HEADER include/ttssimplelib/ttssimplelib.h)

#----------------------------------------------------------------------------------------------------------------------
# ttssimplelib target
#----------------------------------------------------------------------------------------------------------------------
install(TARGETS ttssimplelib
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    PUBLIC_HEADER DESTINATION ${CMAKE_INSTALL_INCLUDEDIR})